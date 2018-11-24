
rm -R crypto-config/*

../bin/cryptogen generate --config=crypto-config.yaml

../bin/configtxgen -profile TFBCOrgOrdererGenesis -outputBlock ./config/genesis.block

../bin/configtxgen -profile TFBCOrgChannel -outputCreateChannelTx ./config/tfbcchannel.tx -channelID tfbcchannel

docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up


## Create channel

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=BankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank.tfbc.com/users/Admin@bank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bank.tfbc.com:7051" cli peer channel create -o orderer.tfbc.com:7050 -c tfbcchannel -f /etc/hyperledger/configtx/tfbcchannel.tx

sleep 5

# Join peer0.bank.tfbc.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=BankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank.tfbc.com/users/Admin@bank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bank.tfbc.com:7051" cli peer channel join -b tfbcchannel.block

# Join peer0.buyer.tfbc.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=BuyerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.tfbc.com/users/Admin@buyer.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.buyer.tfbc.com:7051" cli peer channel join -b tfbcchannel.block

# Join peer0.seller.tfbc.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=SellerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.tfbc.com/users/Admin@seller.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.seller.tfbc.com:7051" cli peer channel join -b tfbcchannel.block
sleep 5

# install chaincode
# Install code on bank peer
docker exec -e "CORE_PEER_LOCALMSPID=BankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank.tfbc.com/users/Admin@bank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bank.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang

# Install code on buyer peer
docker exec -e "CORE_PEER_LOCALMSPID=BuyerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/buyer.tfbc.com/users/Admin@buyer.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.buyer.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang

# Install code on seller peer
docker exec -e "CORE_PEER_LOCALMSPID=SellerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/seller.tfbc.com/users/Admin@seller.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.seller.tfbc.com:7051" cli peer chaincode install -n tfbccc -v 1.0 -p github.com/tfbc/go -l golang

sleep 5


docker exec -e "CORE_PEER_LOCALMSPID=BankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank.tfbc.com/users/Admin@bank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bank.tfbc.com:7051" cli peer chaincode instantiate -o orderer.tfbc.com:7050 -C tfbcchannel -n tfbccc -l golang -v 1.0 -c '{"Args":[""]}' -P "OR ('BankMSP.member','BuyerMSP.member','SellerMSP.member')"
sleep 10
docker exec -e "CORE_PEER_LOCALMSPID=BankMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank.tfbc.com/users/Admin@bank.tfbc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bank.tfbc.com:7051" cli peer chaincode invoke -o orderer.example.com:7050 -C tfbcchannel -n tfbccc -c '{"function":"requestTrade","Args":["A001", "100", "Goods"]}'
