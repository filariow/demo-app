package awsconfig

type Config struct {
	AccessKeyID     string `sbc-key:"access_key_id"`
	SecretAccessKey string `sbc-key:"secret_access_key"`
	SessionToken    string `sbc-key:"session_token"`
	Source          string `sbc-key:"source"`
}
