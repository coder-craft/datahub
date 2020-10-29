package device

var DeviceMgr = &DeviceManager{}

type DeviceManager struct {
}

func (dm *DeviceManager) Name() string{
	return "DeviceManager"
}
func (dm *DeviceManager) Init() bool{
	return true
}
func (dm *DeviceManager) Update() bool{
	return true
}
func (dm *DeviceManager) End() bool{
	return true
}