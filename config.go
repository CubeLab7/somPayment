package somPayment

type Config struct {
	IdleConnTimeoutSec int
	RequestTimeoutSec  int
	Login              string
	Pass               string
	Key                string
	URI                string
}
