package main

import (
	"crypto/x509"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
)

func getTxCreatorInfo(stub shim.ChaincodeStubInterface) (string, string, string, error) {
	var mspID string
	var err error
	var cert *x509.Certificate

	mspID, err = cid.GetMSPID(stub)
	if err != nil {
		fmt.Printf("Error getting MSP identity: %s\n", err.Error())
		return "", "", "", err
	}

	cert, err = cid.GetX509Certificate(stub)

	if err != nil {
		fmt.Printf("Error getting client certificate: %s\n", err.Error())
		return "", "", "", err
	}

	return mspID, cert.Issuer.CommonName, cert.Subject.CommonName, nil
}

// For now, just hardcode an ACL
// We will support attribute checks in an upgrade

func authenticateRealEstateOrg(mspID string, certCN string) bool {

	//return true
	return (mspID == "RealEstateOrgMSP") && (certCN == "ca.realestate-org")
}

func authenticateRegulatorOrg(mspID string, certCN string) bool {
	//return true
	return (mspID == "RegulatorOrgMSP") && (certCN == "ca.regulator-org")
}

func authenticateTraderOrg(mspID string, certCN string) bool {
	//return true

	return (mspID == "TraderOrgMSP") && (certCN == "ca.trader-org")
}

func authenticateShareholderOrg(mspID string, certCN string) bool {
	return authenticateRealEstateOrg(mspID, certCN) || authenticateTraderOrg(mspID, certCN)
}
