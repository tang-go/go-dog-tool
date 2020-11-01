package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.39", nil, nil)
	if err != nil {
		panic(err.Error())
	}
	path := "./dockerfile1"
	var files []*os.File
	rd, err := ioutil.ReadDir(path)
	for _, fi := range rd {
		fmt.Println(fi.Name())
		files = append(files, fi)
	}
	Compress(fi, "")
	if e := Tar(files, "Dockerfile.tar"); e != nil {
		panic(e)
	}
	BuildImage(cli, "./Dockerfile.tar", "", "test/service/docker:v1")
}

//Compress 压缩 使用gzip压缩成tar.gz
func Compress(files []*os.FileInfo, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := compress(file, "", tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.FileInfo, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//BuildImage 编译镜像
func BuildImage(cli *client.Client, tarFile, project, imageName string) error {
	dockerBuildContext, err := os.Open(tarFile)
	if err != nil {
		panic(err.Error())
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Labels: map[string]string{
			project: "project",
		},
	}
	output, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		panic(err.Error())
	}
	scanner := bufio.NewScanner(output.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

//Tar 打包
func Tar(src []string, dst string) error {
	// 创建tar文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	// 通过fw创建一个tar.Writer
	tw := tar.NewWriter(fw)
	// 如果关闭失败会造成tar包不完整
	defer func() {
		if err := tw.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, fileName := range src {
		// 获取文件或目录信息
		fi, er := os.Stat(fileName)
		if er != nil {
			return er
		}

		// 获取要打包的文件或目录的所在位置和名称
		// srcBase, srcRelative := path.Split(filepath.Clean(src))

		srcBase := filepath.Dir(filepath.Clean(fileName))
		srcRelative := filepath.Base(filepath.Clean(fileName))

		// 开始打包
		if fi.IsDir() {
			tarDir2(srcBase, srcRelative, tw, fi)
		} else {
			tarFile2(srcBase, srcRelative, tw, fi)
		}

		// fi, err := os.Stat(fileName)
		// if err != nil {
		// 	log.Println(err)
		// 	continue
		// }
		// hdr, err := tar.FileInfoHeader(fi, "")
		// // 将tar的文件信息hdr写入到tw
		// err = tw.WriteHeader(hdr)
		// if err != nil {
		// 	return err
		// }

		// // 将文件数据写入
		// fs, err := os.Open(fileName)
		// if err != nil {
		// 	return err
		// }
		// if _, err = io.Copy(tw, fs); err != nil {
		// 	return err
		// }
		// fs.Close()
	}
	return nil
}

// 将文件或目录打包成 .tar 文件
// src 是要打包的文件或目录的路径
// dstTar 是要生成的 .tar 文件的路径
// failIfExist 标记如果 dstTar 文件存在，是否放弃打包，如果否，则会覆盖已存在的文件
func Tar2(src string, dstTar string, failIfExist bool) (err error) {
	// 清理路径字符串
	src = path.Clean(src)

	// 判断要打包的文件或目录是否存在
	// if !Exists(src) {
	// 	return errors.New("要打包的文件或目录不存在：" + src)
	// }

	// // 判断目标文件是否存在
	// if FileExists(dstTar) {
	// 	if failIfExist { // 不覆盖已存在的文件
	// 		return errors.New("目标文件已经存在：" + dstTar)
	// 	} else { // 覆盖已存在的文件
	// 		if er := os.Remove(dstTar); er != nil {
	// 			return er
	// 		}
	// 	}
	// }

	// 创建空的目标文件
	fw, er := os.Create(dstTar)
	if er != nil {
		return er
	}
	defer fw.Close()

	// 创建 tar.Writer，执行打包操作
	tw := tar.NewWriter(fw)
	defer func() {
		// 这里要判断 tw 是否关闭成功，如果关闭失败，则 .tar 文件可能不完整
		if er := tw.Close(); er != nil {
			err = er
		}
	}()

	// 获取文件或目录信息
	fi, er := os.Stat(src)
	if er != nil {
		return er
	}

	// 获取要打包的文件或目录的所在位置和名称
	// srcBase, srcRelative := path.Split(filepath.Clean(src))

	srcBase := filepath.Dir(filepath.Clean(src))
	srcRelative := filepath.Base(filepath.Clean(src))

	// 开始打包
	if fi.IsDir() {
		tarDir2(srcBase, srcRelative, tw, fi)
	} else {
		tarFile2(srcBase, srcRelative, tw, fi)
	}

	return nil
}

// 因为要执行遍历操作，所以要单独创建一个函数
func tarDir2(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	// 获取完整路径
	srcFull := srcBase + srcRelative

	// 在结尾添加 "/"
	last := len(srcRelative) - 1
	if srcRelative[last] != os.PathSeparator {
		srcRelative += string(os.PathSeparator)
	}

	// 获取 srcFull 下的文件或子目录列表
	fis, er := ioutil.ReadDir(srcFull)
	if er != nil {
		return er
	}

	// 开始遍历
	for _, fi := range fis {
		if fi.IsDir() {
			tarDir2(srcBase, srcRelative+fi.Name(), tw, fi)
		} else {
			tarFile2(srcBase, srcRelative+fi.Name(), tw, fi)
		}
	}

	// 写入目录信息
	if len(srcRelative) > 0 {
		hdr, er := tar.FileInfoHeader(fi, "")
		if er != nil {
			return er
		}
		hdr.Name = srcRelative

		hdr.Format = tar.FormatGNU

		if er = tw.WriteHeader(hdr); er != nil {
			return er
		}
	}

	return nil
}

// 因为要在 defer 中关闭文件，所以要单独创建一个函数
func tarFile2(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	// 获取完整路径
	srcFull := srcBase + srcRelative

	// 写入文件信息
	hdr, er := tar.FileInfoHeader(fi, "")
	if er != nil {
		return er
	}
	hdr.Name = srcRelative
	hdr.Format = tar.FormatGNU

	if er = tw.WriteHeader(hdr); er != nil {
		return er
	}

	// 打开要打包的文件，准备读取
	fr, er := os.Open(srcFull)
	if er != nil {
		return er
	}
	defer fr.Close()

	// 将文件数据写入 tw 中
	if _, er = io.Copy(tw, fr); er != nil {
		return er
	}
	return nil
}
