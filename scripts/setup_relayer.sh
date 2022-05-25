# We pull and install the relayer
cd ../../
git clone https://github.com/SimplyVC/relayer
cd relayer && make install
cd ../ && rm -r relayer

rly config init

rly chains add --file '../chains/source_chain.json'
rly chains add --file '../chains/target_chain.json'

rly keys restore source-chain default "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"
rly keys restore target-chain default "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

echo "Verify that relayer is funded for both chains"

rly q balance source-chain
rly q balance target-chain

echo "Relayer setup with keys!"