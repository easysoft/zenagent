package repo

import (
	v1 "github.com/easysoft/zv/cmd/server/router/v1"
	"github.com/easysoft/zv/internal/comm/const"
	"github.com/easysoft/zv/internal/server/model"
	"gorm.io/gorm"
	"time"
)

type VmRepo struct {
	BaseRepo
	DB *gorm.DB `inject:""`
}

func NewVmRepo() *VmRepo {
	return &VmRepo{}
}

func (r VmRepo) Register(vm v1.VmRegisterReq) (po model.Vm, statusChanged bool, err error) {
	err = r.DB.Model(&model.Vm{}).
		Where("mac_address=?", vm.MacAddress).First(&po).Error
	srcStatus := po.Status

	// update ip if empty
	r.DB.Model(&model.Vm{}).
		Where("mac_address=? AND (node_ip IS NULL OR node_ip = '')", vm.MacAddress).
		Updates(map[string]interface{}{"node_ip": vm.Ip})

	// just update status by mac for exist vm
	r.DB.Model(&model.Vm{}).
		Where("mac_address=?", vm.MacAddress).
		Updates(
			map[string]interface{}{
				"status":             vm.Status,
				"work_dir":           vm.WorkDir,
				"node_port":          vm.Port,
				"last_register_time": time.Now()})

	err = r.DB.Model(&model.Vm{}).
		Where("mac_address=?", vm.MacAddress).First(&po).Error
	dictStatus := po.Status

	statusChanged = srcStatus != dictStatus

	return
}

func (r VmRepo) GetById(id uint) (vm model.Vm) {
	r.DB.Model(&model.Vm{}).Where("ID=?", id).First(&vm)
	return
}
func (r VmRepo) GetByMac(mac string) (vm model.Vm) {
	r.DB.Model(&model.Vm{}).Where("mac=?", mac).First(&vm)
	return
}

func (r VmRepo) Save(po *model.Vm) {
	r.DB.Model(&model.Vm{}).Omit("").Create(po)
	return
}

func (r VmRepo) UpdateVmName(vm model.Vm) {
	r.DB.Model(&model.Vm{}).Where("id=?", vm.ID).Update("name", vm.Name)
}
func (r VmRepo) UpdateVmCloudInst(vm model.Vm) {
	r.DB.Model(&model.Vm{}).Where("id=?", vm.ID).
		Updates(map[string]interface{}{
			"cloud_inst_id": vm.CloudInstId,
			"node_ip":       vm.NodeIp,
			"mac_address":   vm.MacAddress,
		})
}

func (r VmRepo) Launch(vncAddress, imagePath, backingPath string, id uint) {
	r.DB.Model(&model.Vm{}).Where("id=?", id).
		Updates(
			map[string]interface{}{
				"status":       consts.VmLaunch,
				"vnc_address":  vncAddress,
				"image_path":   imagePath,
				"backing_path": backingPath,
			})

	return
}

func (r VmRepo) UpdateStatusByIds(ids []uint, status consts.VmStatus) {
	db := r.DB.Model(&model.Vm{}).Where("id IN (?)", ids)

	if status == consts.VmRunning {
		db.Where("status <> ?", consts.VmReady) // not to update active vm status
	}

	db.Updates(map[string]interface{}{"status": status})
}

func (r VmRepo) UpdateStatusByNames(vms []string, status consts.VmStatus) {
	db := r.DB.Model(&model.Vm{}).Where("name IN (?)", vms)

	if status == consts.VmRunning {
		db.Where("status <> ?", consts.VmReady) // not to update active vm status
	}

	db.Updates(map[string]interface{}{"status": status})
}

func (r VmRepo) UpdateStatusByCloudInstId(ids []string, status consts.VmStatus) {
	db := r.DB.Model(&model.Vm{}).Where("cloud_inst_id IN (?)", ids)

	if status == consts.VmRunning {
		db.Where("status <> ?", consts.VmReady) // not to update active vm status
	}

	db.Updates(map[string]interface{}{"status": status})
}

func (r VmRepo) DestroyMissedVmsStatus(vms []string, hostId uint) {
	db := r.DB.Model(&model.Vm{}).
		Where("host_id = ? AND status <> ? "+
			" AND strftime('%s','now') - strftime('%s',created_at) > ?",
			hostId, consts.VmDestroy, consts.AgentCheckInterval)

	if len(vms) > 0 {
		db.Where("name NOT IN (?)", vms)
	}

	db.Updates(map[string]interface{}{"status": consts.VmDestroy})
}

func (r VmRepo) DestroyTimeoutVms() {
	db := r.DB.Model(&model.Vm{}).
		Where("status <> ? AND "+
			" strftime('%s','now') - strftime('%s',created_at) > ?",
			consts.VmDestroy, consts.WaitVmLifecycleTimeout)

	db.Updates(map[string]interface{}{"status": consts.VmDestroy})
}

func (r VmRepo) FailToCreate(id uint, msg string) {
	r.DB.Model(&model.Vm{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"status": consts.VmFailCreate, "desc": msg})
}

func (r VmRepo) Distroy(id uint) {
	r.DB.Model(&model.Vm{}).Where("id=?", id).Update("status", consts.VmDestroy)
}
func (r VmRepo) FailToDistroy(id uint) {
	r.DB.Model(&model.Vm{}).Where("id=?", id).Update("status", consts.VmDestroyFail)
}
