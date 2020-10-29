package sign

import (
	"archive/zip"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Sign() {
	// 压缩包
	var src = viper.GetString("in")
	// 解压后保存的位置，为空表示当前目录
	var dst = viper.GetString("out")

	if viper.GetBool("v1-signing-enabled") {
		ChannelV1(src, dst)
		return
	}

	if err := UnZip(dst, src); err != nil {
		log.Fatalln(err)
	}
}

func UnZip(dst, src string) (err error) {

	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
	zr, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	defer zr.Close()

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		path := filepath.Join(dst, file.Name)

		_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			// 因为是目录，跳过当前循环，因为后面都是文件的处理
			continue
		}

		// 创建要写出的文件对应的 Write
		if err = func() error {

			// 获取到 Reader
			fr, err := file.Open()
			if err != nil {
				return err
			}
			defer fr.Close()

			fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer fw.Close()

			n, err := io.Copy(fw, fr)
			if err != nil {
				return err
			}

			// 将解压的结果输出
			fmt.Printf("成功解压 %s ，共写入了 %d 个字符的数据\n", path, n)

			// 因为是在循环中，无法使用 defer ，直接放在最后
			// 不过这样也有问题，当出现 err 的时候就不会执行这个了，
			// 可以把它单独放在一个函数中，这里是个实验，就这样了

			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}
