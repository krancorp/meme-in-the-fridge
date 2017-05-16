# expects target ip address:port as argument, example : "127.0.0.1:8080"
./fridge-sensor $1 Pfungstaedter :9091 &
./fridge-sensor $1 Krombacher :9092 &
./fridge-sensor $1 Grohe :9093 &
./fridge-sensor $1 Salz :9094 &
