package main

import (
        "fmt"
        "log"
        "net"
        "net/http"
        "os"

        "github.com/go-openapi/swag"
        "github.com/s0dimm/aws-s3-proxy/internal/config"
        "github.com/s0dimm/aws-s3-proxy/internal/controllers"
        common "github.com/s0dimm/aws-s3-proxy/internal/http"
        "github.com/s0dimm/aws-s3-proxy/internal/service"
)

var (
        ver    = "dev"
        commit string
        date   string
)

func main() {
        validateAwsConfigurations()

        http.Handle("/", common.WrapHandler(controllers.AwsS3))

        http.HandleFunc("/--version", func(w http.ResponseWriter, r *http.Request) {
                if len(commit) > 0 && len(date) > 0 {
                        fmt.Fprintf(w, "%s-%s (built at %s)\n", ver, commit, date)
                        return
                }
                fmt.Fprintln(w, ver)
        })

        addr := net.JoinHostPort(config.Config.Host, config.Config.Port)
        log.Printf("[service] listening on %s", addr)

        if len(config.Config.SslCert) > 0 && len(config.Config.SslKey) > 0 {
                log.Fatal(http.ListenAndServeTLS(
                        addr, config.Config.SslCert, config.Config.SslKey, nil,
                ))
        } else {
                log.Fatal(http.ListenAndServe(addr, nil))
        }
}

func validateAwsConfigurations() {
        if len(os.Getenv("AWS_ACCESS_KEY_ID")) == 0 {
                log.Print("Not defined environment variable: AWS_ACCESS_KEY_ID")
        }
        if len(os.Getenv("AWS_SECRET_ACCESS_KEY")) == 0 {
                log.Print("Not defined environment variable: AWS_SECRET_ACCESS_KEY")
        }
        if len(os.Getenv("AWS_S3_BUCKET")) == 0 {
                log.Fatal("Missing required environment variable: AWS_S3_BUCKET")
        }

        if swag.IsZero(config.Config.AwsRegion) {
                config.Config.AwsRegion = "us-east-1"
                if region, err := service.GuessBucketRegion(config.Config.S3Bucket); err == nil {
                        config.Config.AwsRegion = region
                }
        }

        if len(config.Config.S3Endpoint) == 0 {
                log.Fatal("S3 Endpoint is not defined in the configuration")
        }
        if len(config.Config.S3AccessKeyID) == 0 || len(config.Config.S3AccessSecretKey) == 0 {
                log.Fatal("AWS Access credentials (AccessKeyID and SecretKey) are not defined in the configuration")
        }

        log.Printf("[config] S3 Endpoint: %s", config.Config.S3Endpoint)
        log.Printf("[config] AWS Region: %s", config.Config.AwsRegion)
        log.Printf("[config] S3 Bucket: %s", config.Config.S3Bucket)
}
