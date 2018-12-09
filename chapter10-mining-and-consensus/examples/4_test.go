package examples_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CoinbaseTransaction() ([]byte, error) {
	const txHash = "d5ada064c6417ca25c4308bd158c34b77e1c0eca2a73cda16c737e7424afba2f"

	resp, err := http.Get("https://blockchain.info/rawtx/" + txHash)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	block, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	var out bytes.Buffer
	err = json.Indent(&out, block, "", "  ")

	return out.Bytes(), err
}

func ExampleCoinbaseTransaction() {
	coinbaseTx, err := CoinbaseTransaction()
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(string(coinbaseTx))

	// Output:
	// {
	//   "ver": 1,
	//   "inputs": [
	//     {
	//       "sequence": 4294967295,
	//       "witness": "",
	//       "script": "03443b0403858402062f503253482f"
	//     }
	//   ],
	//   "weight": 440,
	//   "block_height": 277316,
	//   "relayed_by": "98.117.76.152",
	//   "out": [
	//     {
	//       "spent": true,
	//       "spending_outpoints": [
	//         {
	//           "tx_index": 47893918,
	//           "n": 0
	//         }
	//       ],
	//       "tx_index": 47855746,
	//       "type": 0,
	//       "addr": "1MxTkeEP2PmHSMze5tUZ1hAV3YTKu2Gh1N",
	//       "value": 2509094928,
	//       "n": 0,
	//       "script": "2102aa970c592640d19de03ff6f329d6fd2eecb023263b9ba5d1b81c29b523da8b21ac"
	//     }
	//   ],
	//   "lock_time": 0,
	//   "size": 110,
	//   "double_spend": false,
	//   "block_index": 339688,
	//   "time": 1388185796,
	//   "tx_index": 47855746,
	//   "vin_sz": 1,
	//   "hash": "d5ada064c6417ca25c4308bd158c34b77e1c0eca2a73cda16c737e7424afba2f",
	//   "vout_sz": 1
	// }
}
