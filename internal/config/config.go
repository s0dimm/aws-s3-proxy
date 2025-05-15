package config

import (
        "encoding/json"
        "io/ioutil"
        "log"
        "os"
        "time"
)

type SecretS3 struct {
        Endpoint       string `json:"endpoint"`
        Region         string `json:"region"`
        AccessKeyID    string `json:"accessKeyID"`
        AccessSecretKey string `json:"accessSecretKey"`
}

type BucketInfo struct {
        Spec struct {
                SecretS3 SecretS3 `json:"secretS3"`
        } `json:"spec"`
}

// Config represents its configurations
var (
        Config *config
)

func init() {
        Setup()
}

type config struct {
        AwsRegion          string        // AWS_REGION
        AwsAPIEndpoint     string        // AWS_API_ENDPOINT
        S3Bucket           string        // AWS_S3_BUCKET
        S3KeyPrefix        string        // AWS_S3_KEY_PREFIX
        IndexDocument      string        // INDEX_DOCUMENT
        DirectoryListing    bool          // DIRECTORY_LISTINGS
        DirListingFormat    string        // DIRECTORY_LISTINGS_FORMAT
        HTTPCacheControl   string        // HTTP_CACHE_CONTROL (max-age=86400, no-cache ...)
        HTTPExpires        string        // HTTP_EXPIRES (Thu, 01 Dec 1994 16:00:00 GMT ...)
        BasicAuthUser      string        // BASIC_AUTH_USER
        BasicAuthPass      string        // BASIC_AUTH_PASS
        Port               string        // APP_PORT
        Host               string        // APP_HOST
        AccessLog          bool          // ACCESS_LOG
        SslCert            string        // SSL_CERT_PATH
        SslKey             string        // SSL_KEY_PATH
        StripPath          string        // STRIP_PATH
        ContentEncoding    bool          // CONTENT_ENCODING
        CorsAllowOrigin    string        // CORS_ALLOW_ORIGIN
        CorsAllowMethods   string        // CORS_ALLOW_METHODS
        CorsAllowHeaders   string        // CORS_ALLOW_HEADERS
        CorsMaxAge         int64         // CORS_MAX_AGE
        HealthCheckPath    string        // HEALTHCHECK_PATH
        AllPagesInDir      bool          // GET_ALL_PAGES_IN_DIR
        MaxIdleConns       int           // MAX_IDLE_CONNECTIONS
        IdleConnTimeout    time.Duration // IDLE_CONNECTION_TIMEOUT
        DisableCompression bool          // DISABLE_COMPRESSION
        InsecureTLS        bool          // Disables TLS validation on request endpoints.
        JwtSecretKey       string        // JWT_SECRET_KEY

        // New fields for AWS S3 configuration
        S3Endpoint       string
        S3AccessKeyID    string
        S3AccessSecretKey string
}

// Setup configurations with environment variables
func Setup() {
        // Step 1: Read and parse the /data/cosi/BucketInfo file
        bucketInfo, err := ioutil.ReadFile("/data/cosi/BucketInfo")
        if err != nil {
                log.Fatalf("Failed to read BucketInfo file: %v", err)
        }

        var bucketData BucketInfo
        if err := json.Unmarshal(bucketInfo, &bucketData); err != nil {
                log.Fatalf("Failed to parse BucketInfo file: %v", err)
        }

        // Step 2: Set the values from BucketInfo JSON into config
        Config = &config{
                AwsRegion:          os.Getenv("AWS_REGION"),
                AwsAPIEndpoint:     os.Getenv("AWS_API_ENDPOINT"),
                S3Bucket:           os.Getenv("AWS_S3_BUCKET"),
                S3KeyPrefix:        os.Getenv("AWS_S3_KEY_PREFIX"),
                IndexDocument:      os.Getenv("INDEX_DOCUMENT"),
                DirectoryListing:    false,
                DirListingFormat:    os.Getenv("DIRECTORY_LISTINGS_FORMAT"),
                HTTPCacheControl:   os.Getenv("HTTP_CACHE_CONTROL"),
                HTTPExpires:        os.Getenv("HTTP_EXPIRES"),
                BasicAuthUser:      os.Getenv("BASIC_AUTH_USER"),
                BasicAuthPass:      os.Getenv("BASIC_AUTH_PASS"),
                Port:               os.Getenv("APP_PORT"),
                Host:               os.Getenv("APP_HOST"),
                AccessLog:          false,
                SslCert:            os.Getenv("SSL_CERT_PATH"),
                SslKey:             os.Getenv("SSL_KEY_PATH"),
                StripPath:          os.Getenv("STRIP_PATH"),
                ContentEncoding:    true,
                CorsAllowOrigin:    os.Getenv("CORS_ALLOW_ORIGIN"),
                CorsAllowMethods:   os.Getenv("CORS_ALLOW_METHODS"),
                CorsAllowHeaders:   os.Getenv("CORS_ALLOW_HEADERS"),
                CorsMaxAge:         600,
                HealthCheckPath:    os.Getenv("HEALTHCHECK_PATH"),
                AllPagesInDir:      false,
                MaxIdleConns:       150,
                IdleConnTimeout:    time.Duration(10) * time.Second,
                DisableCompression: true,
                InsecureTLS:        false,
                JwtSecretKey:       os.Getenv("JWT_SECRET_KEY"),

                // Step 3: Assign the S3-related variables from the BucketInfo file
                S3Endpoint:        bucketData.Spec.SecretS3.Endpoint,
                S3AccessKeyID:     bucketData.Spec.SecretS3.AccessKeyID,
                S3AccessSecretKey: bucketData.Spec.SecretS3.AccessSecretKey,
        }

        // Logging the values for verification
        log.Printf("[config] S3 Endpoint: %v", Config.S3Endpoint)
        log.Printf("[config] S3 Access Secret Key: %v", Config.S3AccessSecretKey)
}
