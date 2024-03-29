package abac

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	apValue int
	color   string
	id      string
	owner   string
	size    int
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	_, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	assets := []Asset{
		{id: "asset1", color: "blue", owner: "green", size: 1, apValue: 100},
		{id: "asset2", color: "orange", owner: "bob", size: 10, apValue: 130},
		{id: "asset3", color: "green", owner: "bill", size: 12, apValue: 170},
		{id: "asset4", color: "red", owner: "lucy", size: 11, apValue: 300},
		{id: "asset5", color: "white", owner: "alice", size: 21, apValue: 3400},
		{id: "asset6", color: "brown", owner: "victoria", size: 231, apValue: 1200},
		{id: "asset7", color: "black", owner: "betty", size: 15, apValue: 109},
	}
	fmt.Printf("Trying to populate the ledger...\n")
	for _, asset := range assets {
		asset_json, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.id, asset_json)
		if err != nil {
			return fmt.Errorf("failed to put to the world state. %v", err)
		}
		fmt.Printf("suceessfully put to the world state: ")
		fmt.Printf("id: %s, color: %s, owner:  %s, size: %d, apValue: %d\n", asset.id, asset.color, asset.owner, asset.size, asset.apValue)
	}
	return nil
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	asset_json, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from the world state: %v", err)
	}
	return asset_json != nil, nil
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, owner string, size int, apValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		id:      id,
		color:   color,
		owner:   owner,
		size:    size,
		apValue: apValue,
	}

	asset_json, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, asset_json)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	asset_json, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from the world state: %v", err)
	}
	if asset_json == nil {
		return nil, fmt.Errorf("the asset %s does not exists\n", id)
	}

	var asset Asset
	err = json.Unmarshal(asset_json, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, owner string, size int, apValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exists", id)
	}
	asset := Asset{
		id:      id,
		color:   color,
		owner:   owner,
		size:    size,
		apValue: apValue,
	}
	asset_json, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, asset_json)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exists", id)
	}
	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}
	asset.owner = newOwner
	asset_json, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, asset_json)
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var assets []*Asset

	for resultsIterator.HasNext() {
		var asset Asset
		QR, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(QR.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	return assets, nil
}

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

// func main() {
// 	// 创建并初始化链码实例
// 	myChaincode, err := contractapi.NewChaincode(&SmartContract{})
// 	if err != nil {
// 		log.Panicf("Error creating chaincode:%v", err)
// 	}
// 	// 启动链码
// 	err = myChaincode.Start()
// 	if err != nil {
// 		log.Panicf("Error starting chaincode: %v", err)
// 	}

// }
