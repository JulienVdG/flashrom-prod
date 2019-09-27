#!/bin/sh

cat << EOF
flashrom v1.1-rc1-68-gea0c093 on Linux 4.14.98-v7+ (armv7l)
flashrom is free software, get the source code at https://flashrom.org

Using clock_gettime for delay loops (clk_id: 1, resolution: 1ns).
EOF

case $1 in
	none)
		cat << EOF
No EEPROM/flash device found.
Note: flashrom can never write if the flash chip isn't found automatically.
EOF
		exit 1
		;;
	multi)
		cat << EOF
Found Macronix flash chip "MX25L12805D" (16384 kB, SPI) on linux_spi.
Found Macronix flash chip "MX25L12835F/MX25L12845E/MX25L12865E" (16384 kB, SPI) on linux_spi.
Multiple flash chip definitions match the detected chip(s): "MX25L12805D", "MX25L12835F/MX25L12845E/MX25L12865E"
Please specify which chip definition to use with the -c <chipname> option.
EOF
		exit 1
		;;
	*)
		cat << EOF
Found Macronix flash chip "MX25L25635F" (32768 kB, SPI) on linux_spi.
EOF
esac

# TODO simulate errors

echo -n "Reading old flash chip contents... "
sleep 1
echo "done."
echo -n "Erasing and writing flash chip... "
sleep 2
echo "Erase/write done."
echo -n "Verifying flash... "
sleep 1
echo "VERIFIED."

