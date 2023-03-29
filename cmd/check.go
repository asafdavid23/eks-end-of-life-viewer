/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check for end of life service",
	Long:  `check for end of life service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check called")

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		})

		if err != nil {
			panic(err)
		}

		svc := eks.New(sess)
		input := &eks.ListClustersInput{}

		result, err := svc.ListClusters(input)

		if err != nil {
			errorHandler(err)
		}

		for _, cluster := range result.Clusters {

			input := &eks.DescribeClusterInput{
				Name: cluster,
			}

			result, err := svc.DescribeCluster(input)

			if err != nil {
				errorHandler(err)
			}

			url := "https://endoflife.date/api/amazon-eks/" + *result.Cluster.Version + ".json"

			respData, err := http.Get(url)

			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}

			defer respData.Body.Close()

			body := make([]byte, respData.ContentLength)
			_, err = respData.Body.Read(body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}

			type Product struct {
				EoL string
			}

			var product Product

			err = json.Unmarshal(body, &product)

			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}

			fmt.Printf("Current EKS version for cluster %s is %s and will be end of life on %s, please upgrade to the next version\n", *result.Cluster.Name, *result.Cluster.Version, product.EoL)
		}
	},
}

func errorHandler(err error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
}

func init() {
	rootCmd.AddCommand(checkCmd)
	// checkCmd.Flags().StringP("service", "s", "", "Service name")
	checkCmd.Flags().StringP("version", "v", "1.24", "Version to check")
	// checkCmd.Flags().BoolP("all-versions", "a", false, "Show all versions")
}
