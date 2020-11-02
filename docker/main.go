package main

import (
	"archive/tar"
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.39", nil, nil)
	if err != nil {
		panic(err.Error())
	}
	path := "./dockerfile1"
	createTar(path, "Dockerfile.tar", false)
	BuildImage(cli, "./Dockerfile.tar", "", "test/service/docker:v1")
}

/***生成***/
func createTar(filesource, filetarget string, deleteIfExist bool) error {
	tarfile, err := os.Create(filetarget)
	if err != nil {
		if err == os.ErrExist {
			if err := os.Remove(filetarget); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer tarfile.Close()
	tarwriter := tar.NewWriter(tarfile)

	files, _ := ioutil.ReadDir(filesource)
	for _, f := range files {
		path := filesource + "/" + f.Name()
		fmt.Println("路径:", path)
		sfileInfo, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		if !sfileInfo.IsDir() {
			if err := tarFile(f.Name(), path, sfileInfo, tarwriter); err != nil {
				panic(err)
			}
		} else {
			if err := tarFolder(path, tarwriter); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func tarFile(directory string, filesource string, sfileInfo os.FileInfo, tarwriter *tar.Writer) error {
	fmt.Println("打包", directory, filesource)
	sfile, err := os.Open(filesource)
	if err != nil {
		return err
	}
	defer sfile.Close()
	header, err := tar.FileInfoHeader(sfileInfo, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	header.Name = directory
	err = tarwriter.WriteHeader(header)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if _, err = io.Copy(tarwriter, sfile); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func tarFolder(directory string, tarwriter *tar.Writer) error {
	baseFolder := filepath.Base(directory)
	return filepath.Walk(directory, func(targetpath string, file os.FileInfo, err error) error {
		//read the file failure
		if file == nil {
			return err
		}
		if file.IsDir() {
			// information of file or folder
			header, err := tar.FileInfoHeader(file, "")
			if err != nil {
				return err
			}
			header.Name = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			if err = tarwriter.WriteHeader(header); err != nil {
				return err
			}
			os.Mkdir(strings.TrimPrefix(baseFolder, file.Name()), os.ModeDir)
			return nil
		}
		//baseFolder is the tar file path
		// trim := strings.TrimPrefix(targetpath, directory)
		// fileFolder := filepath.Join(baseFolder, trim)
		fileFolder := baseFolder + "/" + file.Name()
		//fmt.Println("name", file.Name(), "baseFolder=", baseFolder, "targetpath=", targetpath, "directory=", directory, "fileFolder=", fileFolder)
		return tarFile(fileFolder, targetpath, file, tarwriter)
	})
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
