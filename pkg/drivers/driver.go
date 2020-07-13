package drivers

type Driver interface {
	Init(driverConfig map[string]string)
}
