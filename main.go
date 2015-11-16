package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "s3-get"
	app.Usage = "used to pull a single file from s3 and save it locally"
	app.Version = "1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Daniel Baldwin",
		},
	}
	app.Copyright = `
The MIT License (MIT)
Copyright (c) 2015 MasteryConnect
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
	`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "region,r",
			Value: "us-east-1",
			Usage: "aws region name",
		},
		cli.StringFlag{
			Name:  "bucket,b",
			Usage: "s3 bucket",
		},
		cli.StringFlag{
			Name:  "key,k",
			Usage: "s3 key",
		},
	}
	app.Action = func(c *cli.Context) {
		// Vars
		var path string
		svc := s3.New(session.New(), &aws.Config{Region: aws.String(c.String("region"))})

		// Checks
		if len(c.String("bucket")) == 0 {
			log.Fatal("Need an S3 Bucket.")
			cli.ShowAppHelp(c)
		}
		if len(c.String("key")) == 0 {
			log.Fatal("Need an S3 Key.")
			cli.ShowAppHelp(c)
		}
		if len(c.Args()) > 0 {
			path = strings.Join(c.Args(), " ")
		} else {
			log.Fatal("Need a destination path.")
			cli.ShowAppHelp(c)
		}

		resp, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(c.String("bucket")),
			Key:    aws.String(c.String("key")),
		})
		if err != nil {
			log.Fatal(err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = ioutil.WriteFile(path, buf.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}
