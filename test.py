# simple example of reading the mcp3008 analog input channels and printing
# them all out.
# author: tony dicola
# license: public domain
import time

# import spi library (for hardware spi) and mcp3008 library.
import adafruit_gpio.spi as spi
import adafruit_mcp3008

# hardware spi configuration:
spi_port   = 0
spi_device = 1
mcp = adafruit_mcp3008.mcp3008(spi=spi.spidev(spi_port, spi_device))


print('reading mcp3008 values, press ctrl-c to quit...')
# print nice channel column headers.
print('| {0:>4} | {1:>4} | {2:>4} | {3:>4} | {4:>4} | {5:>4} | {6:>4} | {7:>4} |'.format(*range(8)))
print('-' * 57)
# main program loop.
while true:
    # read all the adc channel values in a list.
    values = [0]*8
    for i in range(8):
        # the read_adc function will get the value of the specified channel (0-7).
        values[i] = mcp.read_adc(i)
    # print the adc values.
    print('| {0:>4} | {1:>4} | {2:>4} | {3:>4} | {4:>4} | {5:>4} | {6:>4} | {7:>4} |'.format(*values))
    # pause for half a second.
    time.sleep(0.5)
