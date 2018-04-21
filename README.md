# Components

* [Raspberry Pi 3 Model B](https://www.raspberrypi.org/products/raspberry-pi-3-model-b/)
* [MCP3008](http://ww1.microchip.com/downloads/en/DeviceDoc/21295C.pdf) wired to hardware SPI

# Setup

1. Download [Raspbian Stretch Lite][0] and unzip
1. Insert SD card into host computer ([note for 30GB+ cards][4])
1. Format SD card to FAT32 using native tool or [official tool][1]
1. Local `diskutil list` to find the SD volume
1. Local `diskutil unmountdisk /dev/disk2` with the actual path from previous step
1. Local `sudo dd if=raspbian.img of=/dev/disk2 bs=2m` with actual paths
1. Wait for quite some time...
1. Local `diskutil mountdisk /dev/disk2` to mount the disk
1. Local `cd /Volumes/boot`
1. Touch `/ssh` on the boot partition somehow to [enable SSH][2]
1. Insert SD card, plug in ethernet, then plug in power... booted!
1. SSH using default credentials pi:raspberry and IP discovered via router tool
1. `sudo apt-get install vim`
1. Uncomment `duid` in `/etc/dhcpcd.conf` on the boot partition to [fix IPv6][3]
1. Run `sudo raspi-config`
    1. Select '1 Change User Password' to change default password
    1. Select '2 Network Options' then 'N2 Wi-fi' to connect to Wi-fi
    1. Select '4 Localisation Options' then 'I2 Change Timezone' to update clock
    1. Select '5 Interfacing Options' then 'P4 SPI' to [enable SPI][5]
    1. Select '7 Advanced Settings' then 'A1 Expand Filesystem' to [expand][42]
1. Run `cd ~ && install -d -m 700 ~/.ssh` to initialize the SSH directory
1. On the host, run `cat ~/.ssh/id_rsa.pub | ssh {user}@{host} 'cat >> .ssh/authorized_keys'`
   to copy local keys to Pi
1. Run `sudo reboot` and cross fingers
1. Wait for Pi to reboot, then SSH back in... (reboot again if SSH wont work)
1. [We expect][5] `lsmod` to indicate the SPI Linux Kernel module is loaded
   `$ lsmod | grep spi_`
1. [We expect][7] an interface to be mounted `$ ls -l /dev/spi`
1. TODO: Compile and test logger program

[0]: https://www.raspberrypi.org/downloads/raspbian/
[1]: https://www.sdcard.org/downloads/formatter_4/
[2]: http://blog.smalleycreative.com/linux/setup-a-headless-raspberry-pi-with-raspbian-jessie-on-os-x/
[3]: https://www.raspberrypi.org/forums/viewtopic.php?f=63&t=177624
[4]: https://www.raspberrypi.org/documentation/installation/sdxc_formatting.md
[5]: https://www.raspberrypi-spy.co.uk/2014/08/enabling-the-spi-interface-on-the-raspberry-pi/
[6]: https://gist.github.com/bloudermilk/f7d5033ad2f7e66c22c993d4e3d00c91
[7]: https://learn.adafruit.com/raspberry-pi-analog-to-digital-converters/mcp3008#hardware-spi
[8]: https://github.com/adafruit/Adafruit_Python_MCP3008#installation
[9]: https://learn.adafruit.com/raspberry-pi-analog-to-digital-converters/mcp3008#library-usage
[42]: https://media3.giphy.com/media/xT0xeJpnrWC4XWblEk/giphy.gif
