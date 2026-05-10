package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhangyiming748/lumberjack"

	"avdb/param"
	"avdb/soup"
	"avdb/storage"
)

// 版本变量，将在构建时注入版本信息
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

/*
main函数用cobra实现一个命令
主命令avdb
子命令file
参数-i --import对应 SearchByJavdb函数的src参数
参数-e --export对应 SearchByJavdb函数的dst参数
参数-c --config对应 param.InitConfig函数的configPath参数
且param.InitConfig函数需要在SearchByJavdb之前被调用
*/
func main() {
	var configFile string
	var importFile string
	var exportFile string

	// 创建根命令
	rootCmd := &cobra.Command{
		Use:     "avdb",
		Short:   "AVDB 是一个成人影片数据库工具",
		Long:    `AVDB 是一个用于查询和管理成人影片数据的工具`,
		Version: Version,
	}

	// 添加版本标志
	rootCmd.Flags().BoolP("version", "v", false, "显示版本信息")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		// 如果用户请求版本信息
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			PrintVersion()
			return
		}
		// 如果没有子命令且没有请求版本，显示帮助
		cmd.Help()
	}

	// 添加 version 子命令
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Long:  `显示 AVDB 的版本信息`,
		Run: func(cmd *cobra.Command, args []string) {
			PrintVersion()
		},
	}
	rootCmd.AddCommand(versionCmd)

	// 添加 file 子命令
	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "处理文件导入导出",
		Long:  `使用指定的导入文件和导出文件处理数据`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// 在执行子命令前初始化配置
			if configFile != "" {
				param.InitConfig(configFile)
				storage.SetSqlite() // 初始化数据库
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// 执行文件处理操作
			result, err := soup.SearchByJavdb(importFile, exportFile)
			if err != nil {
				log.Printf("处理文件时出错: %v", err)
				os.Exit(1)
			}
			log.Printf("文件处理完成，结果保存至: %s", result)
		},
	}

	// 为 file 命令添加标志
	fileCmd.Flags().StringVarP(&importFile, "import", "i", "", "导入文件路径 (必需)")
	fileCmd.MarkFlagRequired("import") // 标记为必需参数

	fileCmd.Flags().StringVarP(&exportFile, "export", "e", "export.txt", "导出文件路径 (可选，默认为export.txt)")
	// fileCmd.MarkFlagRequired("export") // 默认值情况下不需要标记为必需参数

	fileCmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件路径 (可选)")

	// 添加子命令到根命令
	rootCmd.AddCommand(fileCmd)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		log.Printf("执行命令时出错: %v", err)
		os.Exit(1)
	}
}

func PrintVersion() {
	if Version != "dev" {
		fmt.Printf("AVDB 版本: %s\n", Version)
		fmt.Printf("构建时间: %s\n", BuildTime)
		fmt.Printf("Git 提交: %s\n", GitCommit)
	} else {
		// 尝试从 runtime/debug 获取版本信息
		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Printf("AVDB 版本: %s (devel)\n", info.Main.Version)
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					fmt.Printf("Git 提交: %s\n", setting.Value)
				} else if setting.Key == "vcs.time" {
					fmt.Printf("构建时间: %s\n", setting.Value)
				}
			}
		} else {
			fmt.Printf("AVDB 版本: %s\n", Version)
		}
	}
}

func SetLog(l string) {
	// 设置全局时区为Asia/Shanghai
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("无法加载时区 Asia/Shanghai: %v", err)
	} else {
		time.Local = location
	}
	// 创建一个用于写入文件的Logger实例
	fileLogger := &lumberjack.Logger{
		Filename:   l,
		MaxSize:    1, // MB
		MaxBackups: 1,
		MaxAge:     28, // days
	}
	err = fileLogger.Rotate()
	if err != nil {
		log.Println("转换新日志文件失败", err)
	}
	consoleLogger := log.New(os.Stdout, "CONSOLE: ", log.LstdFlags)
	log.SetOutput(io.MultiWriter(fileLogger, consoleLogger.Writer()))
	log.SetFlags(log.Ltime | log.Lshortfile)
}
