#!/usr/bin/env python3

# imports
import os
import time

Y = (['yes', 'y', 'YES', 'Y'])
N = (['no', 'n', 'NO', 'N'])

def install():
    os.system('clear')
    print("INSTALLING")
    time.sleep(1)
    os.system('go get -u github.com/json-iterator/go')
    os.system('go get -u github.com/valyala/fasthttp')
    os.system('go get -u github.com/valyala/tcplisten')
    os.system('go get -u github.com/krishpranav/webfr')
    time.sleep(1)
    print("WEBFR has been installed successfully")
    time.sleep(1)
    print("check out: https://github.com/krishpranav/webfr for more details")
    print("""To check that webfr has installed succesfully enter this command:
            ls go/src/github.com/krishpranav
    """)
    time.sleep(1)
    os.system('cd ../; rm -rf webfr')


def main():
    choice = input("You sure to install webfr in your system Y / N: ")
    if choice in Y:
        install()
    elif choice in N:
        print('EXITING')
        os.exit(0)

if __name__ == "__main__":
    print("WEBFR INSTALLER")
    main()