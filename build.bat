cd server
goxc -os=linux -arch="amd64"
goxc -os=windows -arch="386"
cd ../client
goxc -os=linux -arch="amd64"
goxc -os=windows -arch="386"
cd ..