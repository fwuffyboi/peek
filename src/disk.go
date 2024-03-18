package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// DiskUsage create struct
type DiskUsage struct {
	Path        string
	DiskName    string
	UUID        string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

// storageAllDisks Returns the total, free, and used disk space for all disks
func storageAllDisks() (allDiskUsage []DiskUsage, err error) {

	// get fstab so we can see all disks
	config, err := ConfigParser()
	if err != nil {
		log.Fatalf("Error getting config. err: %s", err)
	}

	// if the config is set to show disk usage
	if config.Show.ShowDisk == true { // get disk usage for all disks

		// get fstab so we can see all disks
		fstab, err := os.Open("/etc/fstab")
		if err != nil {
			log.Errorf("Error opening /etc/fstab. err: %s", err)
			return nil, err
		}
		defer fstab.Close() // todo: error handle

		// get the disk usage for all disks
		scanner := bufio.NewScanner(fstab)

		// loop through all the disks in the fstab, and append when data is gathered
		for scanner.Scan() {

			// get the line
			line := scanner.Text()

			// split the line
			fields := strings.Fields(line)
			if len(fields) < 2 || strings.HasPrefix(fields[0], "#") {
				// Skip lines with less than 2 fields and lines that start with '#'
				continue
			}

			// get values from line fields
			mountPoint := fields[1]
			total, free, used, err := getDiskUsage(mountPoint)
			if err != nil {
				log.Errorf("Couldn't get disk usage for mountpoint: %s, error: %s", mountPoint, err)
				return nil, err
			}

			// get the disk name and uuid
			dfCmd := exec.Command("df", "--output=source", mountPoint)
			dfOut, err := dfCmd.Output()
			if err != nil {
				return nil, err
			}

			device := strings.TrimSpace(string(dfOut))
			device = strings.Split(device, "\n")[1] // The first line is the header

			lsblkCmd := exec.Command("lsblk", "-o", "UUID,NAME", "-nl", device)
			lsblkOut, err := lsblkCmd.Output()
			if err != nil {
				return nil, err
			}

			output := strings.TrimSpace(string(lsblkOut))
			split := strings.Split(output, " ")

			if len(split) < 2 {
				return nil, fmt.Errorf("unexpected output from lsblk: %s", output)
			}

			uuid := split[0]
			name := split[1]

			// calculate used percent
			usedPercent := (float64(used) / float64(total)) * 100

			// log
			log.Infof("DiskPath: %s, DiskName: %s, UUID: %s, Total: %d, Free: %d, Used: %d\n", mountPoint, name, uuid, total, free, used)

			// append to list
			allDiskUsage = append(allDiskUsage, DiskUsage{
				Path:        mountPoint,
				DiskName:    name,
				UUID:        uuid,
				Total:       total,
				Free:        free,
				Used:        used,
				UsedPercent: usedPercent,
			})
		}

	} else {
		log.Info("Disk usage is not set to be shown in the config")
		return nil, nil
	}

	return allDiskUsage, nil

}

func getDiskUsage(path string) (total, free, used uint64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		log.Errorf("Error getting disk usage for path: %s. err: %s", path, err)
		return 0, 0, 0, err
	}

	// Total space.
	total = stat.Blocks * uint64(stat.Bsize)

	// Free space.
	free = stat.Bfree * uint64(stat.Bsize)

	// Used space.
	used = total - free

	return total, free, used, nil
}
