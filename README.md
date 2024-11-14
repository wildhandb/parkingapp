# Parking App With Go

Build & Run (Windows)

```
git clone https://github.com/wildhandb/parkingapp.git
cd parkingapp
go build .

parkingapp create_parking_lot 6
parkingapp park KA-01-HH-1234
parkingapp park KA-01-HH-9999
parkingapp park KA-01-BB-0001
parkingapp park KA-01-HH-7777
parkingapp park KA-01-HH-2701
parkingapp park KA-01-HH-3141
parkingapp leave KA-01-HH-3141 4
parkingapp status
parkingapp park KA-01-P-333
parkingapp park DL-12-AA-9999
parkingapp leave KA-01-HH-1234 4
parkingapp leave KA-01-BB-0001 6
parkingapp leave DL-12-AA-9999 2
parkingapp park KA-09-HH-0987
parkingapp park CA-09-IO-1111
parkingapp park KA-09-HH-0123
parkingapp status
```