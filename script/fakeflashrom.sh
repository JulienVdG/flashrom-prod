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

if [ "$1" = "nofile" ]; then
	echo "Error: opening file "this_file_doesnt_exist" failed: No such file or directory"
	exit 0
fi

# TODO improve error simulation (ie handle multiple erase function)
	echo -n "Reading old flash chip contents... "
	sleep 1
	echo "done."
	echo -n "Erasing and writing flash chip... "
case $1 in
	same-content)
		cat << EOF

Warning: Chip content is identical to the requested image.
Erase/write done.
EOF
		exit 0
		;;
	erase-error)
		sleep 1
		echo "FAILED at 0x00331000! Expected=0xff, Found=0x00, failed byte count from 0x00331000-0x00331fff: 0x1000"
		echo "ERASE FAILED!"
		# no more timing simulation, code currently stops at first error
		cat << EOF
Reading current flash chip contents... done. Looking for another erase function.
FAILED at 0x00000000! Expected=0xff, Found=0x00, failed byte count from 0x00000000-0x00000fff: 0x1000
ERASE FAILED!
Reading current flash chip contents... done. Looking for another erase function.
FAILED at 0x00000000! Expected=0xff, Found=0x00, failed byte count from 0x00000000-0x00007fff: 0x8000
ERASE FAILED!
Reading current flash chip contents... done. Looking for another erase function.
FAILED at 0x00000000! Expected=0xff, Found=0x00, failed byte count from 0x00000000-0x00007fff: 0x8000
ERASE FAILED!
Reading current flash chip contents... done. Looking for another erase function.
FAILED at 0x00000000! Expected=0xff, Found=0x00, failed byte count from 0x00000000-0x0000ffff: 0x10000
ERASE FAILED!
Reading current flash chip contents... done. Looking for another erase function.
FAILED at 0x00000000! Expected=0xff, Found=0x00, failed byte count from 0x00000000-0x0000ffff: 0x10000
ERASE FAILED!
Reading current flash chip contents... done. Looking for another erase function.
FAILED at 0x00000000! Expected=0xff, Found=0x00, failed byte count from 0x00000000-0x01ffffff: 0x2000000
ERASE FAILED!
Reading current flash chip contents... done. Looking for another erase function.
FAILED at 0x00000000! Expected=0xff, Found=0x00, failed byte count from 0x00000000-0x01ffffff: 0x2000000
ERASE FAILED!
Reading current flash chip contents... done. No usable erase functions left.
FAILED!
Uh oh. Erase/write failed. Checking if anything has changed.
Reading current flash chip contents... done.
Apparently at least some data has changed.
Your flash chip is in an unknown state.
Please report this on IRC at chat.freenode.net (channel #flashrom) or
mail flashrom@flashrom.org, thanks!
EOF
		exit 0
		;;
	*)
		sleep 2
		echo "Erase/write done."
		echo -n "Verifying flash... "
		sleep 1
		echo "VERIFIED."
esac
		sleep 1
		echo "JVDG"
