// Package public provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package public

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xce28buXb/KsT0ArcFNLIejuMIWBTe2HcjNK86TrFtZBjUzJHE9Qw5ITl2dA0B/Rr9",
	"ev0kBR/z4nD0cJxs7nb3j4Wj4RwenseP50HOQxCxNGMUqBTB5CHgIDJGBeh/vOSAJZyzFBN6aR+o3yNG",
	"JVCp/sRZlpAIS8Lo0W+CUfWbiFaQYvXXXzgsgknwT0fVJEfmqTgyZIPNZtMLYhARJ5miEkyCYip0vwKK",
	"MIr1SLTCAs0BKIo0W3E/6AUrwDFwzeyv4eUqnMZpeLmCRP9xxW7BYUiuMwgmgZCc0KWae9MLLjhn/MmX",
	"p6kK3/L0E1QIGi0YR4KlgFgGXM/UDza94BUT8iWjiydnzCX8wby1TQ+KRbkCVLyKcFTy+ZoIaVQpnpxV",
	"D+1ubqcSccg4CKBSc9uQMEYJERKxhbUmgT7nwNd6CZeA4+9u5FcrIrYzLDKIyIJElmPL6pIICfz388nC",
	"FqxTcsuQY70fs/h3RQ6Hy1yz0+Bx07OzGaRbQXSrzFv/q0n0HBaEgiZHaJZLFGOJyxmOIvXqzYoJWbjF",
	"jM6oVq/5NyICJTin0QpitOAsNaQyjATwO+BIMqSJ6N81IbnCUr02B0KXM4pzyUKgnCWJgT34gtMs0RI1",
	"K7yhOFXAZk2leN4rHlvc4woaSZwGvYDQO6CS8fUNiYNJMAJ8Es+fQxgtnh+Hx89OBiEenZ6Ex/Hz6AXg",
	"xcmL00XQC0Q+L0Vzk2KKl8AtgfHo5Pn4FIc4XpyGwyHEIV48j8Pj0xGG8enJOB4PlGVkXGlBErPDNLh/",
	"aLlIqUA1QG0HHJS0jCgqcd2lDaEEM0cOsyDoudDviMade0pjZZxG7WpIBR19dBbJHCfJGjGarLUmiUAi",
	"zzLGpaOgYFYK3c9FUw8+EZBYzV1bLCJUMmuPgixX0jwoKbkM7KNbP3Od6vbxWR+M7GDFvEazyoK3qW0f",
	"I/JxuukFHD7nhEMcTD41ZdpUdK9hct0rvO4FkkjN1iVjEl0pE1ArqSFFyQab/waRVAKrR0wa5OKYKNI4",
	"eV+z+wVOBLjYdVZBqmA5j6Dp50qCN0aCKVB5AxTPE7VcyfMWrTdrZHhAF4YCOq89r5bWHneln5QrKwGj",
	"5bmd/LimcaEfIMZRTIT+szDjmlEoEhZSibCCaNiGWmXF1pyxBDAN3A3AnftVnmIlUBzriWsP3cmalqjR",
	"mwiUrhG+Bx2dWeXUSPQPBpUriyP1edHLnHOgskCTQuJtSKF5qoy7VMl1neMasrd4stp2udHIph6VexlQ",
	"ydc7vatL8y1H85hcp1s5oYLHtX4PpzLrafnWdpuwb+ldcYjx8BlgCPHp6dwC2iiaN3fFnrOJp+sDtvFM",
	"h0sRvlHxbTD5dN0LooSoxTDNkXJWrUCcpMUE0/dn/Ytfz968f32hIFBHIcX79p/DvoeLD/rRqP3oevOH",
	"wwc8Z3kthNzm7L4NcUolcIqTYgPs4ntPA+ma2h84nVdBU2Myj0r/IQCsfLZXYjDN8BbU0xucCiLNz3PY",
	"ruRDsc/K3BtymGDkW6Di1KDAAcD4zoCDnkUFsDYrDlw3jvBNBNyXFp2hORZwclwlsTq1UuaBbZRXpNzR",
	"Spmj+kMRXujAWvQQoVGSx4Qu9WCbDEW4uRV/mlGEHtT/EJoFlES3SsSzYIJmwZv/PH/35mz6tgAzNH1/",
	"hl6ezYJeMZ4IkQM3o1++/ellNT86y+WKcSLXvXc/uYRqFESuZf1VJIATnNzQPJ0XvAxrjymTN3c4IfHN",
	"HBaM27WNBqNxOBiG4+HVcDwZjSfjk//yvoQXsqDqf6d4JYPUDAvVfz9f/DJ9i15eXF5N/zZ9eXZ1oX+d",
	"zeib6fSi3+/PZlT/cvH23DfKsrKZ0Wu1C0lIxd7eqUSoPdQYM+Ycr4NyjyossWVwr6016WEIC8EigiXE",
	"yCZDSvmVH5cstcMh77TdKegt8DlwJuzEOhV1ZyfCmb4yYNcwvHmWtv0ta3a3Qos71mmEL5o9TCdma2+L",
	"x8G/Eg4qnhsCdJXYiWwFpk0z/M6A2TZw0xZzGMBdlqU1XIcdpMSCCNUWE2G9GINTRRA8LdFQ5/39FiQa",
	"TFF/VUo+ABh8+i9grUm0A968BJpw0KQT9IIF4ymWwSSIsYRQkhS2UzFI9EgyGaTNN3cCzm686fCaCleb",
	"Mw47qhnasp5Oc27dwVhGTZ9tzXikbCTmrqfit7sWUas2FF6yzY2skx/mSOYlZBDDFHOoyRSrOpTHZ2pZ",
	"VoRvhJ3Z5FWLz7GiTDJsfvflGqsoqgdaTQLqoanpOg8SZirKwSSAXK1MCfeWUOJS4CuRmhRtMH5x8gyP",
	"onAej45NBP7i5HRX4bK2pod90hCz5PrQLXlWy3I7hLHXzB5Z7fVeJcp25Fc8QyQGqnwHOPrnhN0DDyMs",
	"AJ2//YASPIfkX5QppPjLa6BLuQomJ2MfXjQVtBd3WyqT1dt7qXaXU1eatlrsUodP1O7iuhmvCXxff/9Q",
	"0Gx5vG4xthX3C1DgOEGgO5C22dLsJoDEJAkmwc84RkoKIOQEVWiBzAAxQSnj4HmgAg/lVX95+DW8fBVO",
	"336Y/vLq6kN4efHvHy8+XIXT840SgsQyF8EkOB74mgKWlMfwzCOIEXzJEkyNEZaNsiIUY1Gkc9IIisp5",
	"xtk8gbSHoL/sI4wkxxHMcXTbjNXOqJWNpgAxwhJFLIb+kk2Go7G3CODL/s9QTsnnHOr+UYZoGeaSRHmC",
	"eTejTbZuf8txLBY5Fqt7OH1xev/it4V3g7Ni9UWxr66u3iMzQC8J2f7b3GTDhrNSSl9U8CSUAATCyMyA",
	"7nCSQx/p4uiHV+8+vj5Hc83xHYndtodSbMPzx/pfJM3TYPLsxYtekBJaf0aoeTYcDHQQGb+jybpAa7tS",
	"QiUsbZTqz+7P0CpPMQ3Lko5YMS4faS12qQsCSfy///0/wggARZiqdSt3Tcjf3YXXHGd3vBBXvlAsqVea",
	"fycMGO/ucnuxy+/LNjMHmXMKMZqvTVMmFSRO0dn7aRMUwJL99Dvhw3ULIaBjoeZohRGIKCyzzgWhJoQt",
	"Knyy0LVAQCOWUwnK7e9XJAGUAVejleljWm8ed9cnSGxS7C3LqxcG9DLNG8eDQe2JkZZ58lhRm8xfZeg6",
	"fN4rHzSmtdH+ODUvDDvdsUoRt9jptuJVcZ6krchL52CEXj3K8DphOK567ysmZMTo4miidmYENM4YodKJ",
	"P224qQYP+61gq7ulHY1xdDx/MQxHz6MoPB6NohCPTsfh8QjjeAHPxwO8aO9eRaTXhuBFniToc44TtSHE",
	"jda20+Bt9P+Rt/MfzIJ0rcb39+l1H9TIJZWEfT1dl419xNTRtPXbTWkWW0zGOQ/kXQd3T4PY8xiSKejj",
	"BO7AaatX4HC/AjqjxKcF5JztkPgWaHWww9fafvwxjXrp2xZfdiTURT69JZ2mgdMQauW5Tk+oytR2N36c",
	"zsSuBkM7njK1bG8Z2vtCs37XUVuzS3nYuy64aRnfTnvdBnW1o2znWB5ardcnYbA0h9ggdo43if5juprf",
	"oz95SFfwMa270Nu6s2Bq8VWXKna373Z12HbZsa+7tW8n3qxlwQnQOFk7Tfkdy3jyBlUVhlYsdoairllv",
	"t/w2aD+yoJvhJaG6+s5B5In0H/F0EFi73acfyj2ue0FC6K2JGwhXCBUc4YwckTg9uhse2YX8a0JSIn8a",
	"Dmb5YDA6YYuFAPmTmjnBh70zVC9R+HL4SxmHO8J0UH4AhyqSBAN3OroOJkPPjmEBsY159rhoPXdewt7N",
	"Ddc4Pb2fUvrb6Lw31kYYfa2H1xa132tv1GjXUzWJnll8wcg+fubzMZfDx3qWwhuK78jSBEGaqQqI1AZU",
	"ep4bY/9BrdeJ7c0q29nKAkwBwdbD9ThtrMVmZFMY7z5kxLCTqBq2N00jpZ001bC9aVZC3Em3GLon7c7Q",
	"yjXsrbb/Bg4OrNQrjePcjZ3F3IRxuxkWx5TPpkT9NegFxmSCyaDdKDDDW4fvmMQJ4hAxHouyH6kSjKhd",
	"YxjWW2+EypPjwFcTs/y4U73V/SSlBI2aKANe4mibSLGS9oEmIknJclkp66Dj1tC1EAoOy0mud6v9Dfij",
	"iuZFiJ1ax0nybqF3f5fSflcNrvc5WPjYMzz1yxLffyVm9uqa1VcvSOdcUc6JXOtAz/jBl5CvQlMOl2sf",
	"gmhrQRjp1iiyt7ns8F5AdAVH3/ZSe4CJtNxBxe6ekX+DtbkfonL64uIJ1r3f1qWRn3F0CzRGZ++nen22",
	"DKryRtIQg3kQzs34oBfcAReGxrA/6Ov9gmVAcUaCSTDuD/ojtfVgudISMDdFwhUT8uihoyCzOXpYfI7p",
	"RmfAmOMUJJSl1z0P3NewQ897hjimMUtRnpMYLXUtWNavokjgKaE4QfM1+qsa9Ff1mi522/sDpxE8wyfh",
	"fPx8YNto4/HAPZVvxB1M9IorJXU3vCqEMJF39zW93oHltALNdX3MFQdF9p9IFw1N6mBGNpbtqxd2L9O2",
	"Br9iTYUhI2Pm/X2t/iumvLSFVVvxi4kaOs/1sSaOI0KXO7gw909EaAmFU93R6GTguhdkzBfsvMRJsut2",
	"lFFS0X4pynHVIYTq6FNZqp+qJK68sWFlBUL+zOL1k11Hq90I2TQ3PqUP/UPtTu1ocOzPdeztGKcSbFeb",
	"4hj0xbrjwaCLn3Kao2bhXb91fPBbNRjX+OMA+KdrpU2Rpynm60LKhY7uUmTOtjQu/ng0JfFS6LJFZM6q",
	"63riUecRwKlEQiqyJnaNQAisz6G4XZ3qOMoKS3RPkgQlIC1IqrfToBcsfXHO3yx0mJ6gRIwvMSV/18R7",
	"psJQHG8tbpKqOWYUc7CVYjOLteaIUcES6MdMtiyzmdVtxfusjIiQDZ6aqHYOC5wnElXhmwWygQmRAY20",
	"i6uATSgjGw5QmUTbocNBiW36cmzl7OWULdeuRXz7xp0O58Nn7vMaR882vXJp9rc92C3CzC3cXrfccg+/",
	"8t11/pF9Up8eLazU5hh1e645YBH1KRfcYYrOPlVutlXX+//PtqUP5FRoU+RxHuzw7Ev1u4LfaGuqT+G5",
	"L20vipjIXe20c6h/XWHXTjbcbcDeT0j80PtYU6X9Th+pXjl3EqfaFnb0oOLpjbGcBKTvso7+XR+q/EKE",
	"VBuItaYyGtrPngyhmj3tijoui/Mf+jsb9gsEegnoHgtkGI5/7LDDSm+XspwjrcURoFKAAqnQoa0C3SrG",
	"ca/8ggC3YkFEJxbeEKKAJ3NDuwpNvJ91UPsjy/kO3VafqAges3F5vnDxIytVsbuHB/65S+1OknXOb46T",
	"mKOIaycCr4diuKNSUNYFUMTSFNP4KesDivBBgm6YyrumB7ve5cHio+KTJZ7YX5ftzXNDVQGAdf3mFaO6",
	"Y/eDw2zxu9reI0zsK+Z9RwFJkgLKBSDJbsFCqP1IjP5URlNw2ySRhpc1XdjvOX0Vf9qwy5s+vhtb2hG0",
	"TdZsfHA6PhlH8XE4B2y/SzE/HS+cOwJPZN8tpkmGw1UUIbxUQUBRdNwluv8oxx3kWlnu2dE+CmP/zKve",
	"moMQ2yVwtr0ZxdW94K46jVPS/zYRsTOJ7ztS7o5N624vGDKHwsyimyWrSkv7RM977dvezz19xd49/jF3",
	"fKcJ4S8NNUMECyltOEZyxdm9TyeezcCo0WPxlU07hvsn8v+JrE+IrDuxE+0BnY0O4rcBzsYUPzhser8+",
	"9wcGzeK0RbIuQ62yIbYbSi2ZQ4F0xYQM9eH7qmXpftnFHMdX0YJOIoregPNFsNRqfmfDkzg5XftcybdG",
	"14/mWhfvSO30tyy+W263332CkpGv7Fd2FkHPtpyqd7Gs3tpqX25oglrtGsBhgFbfb2JIWZnzV/vOt74U",
	"Un1ect/Psfo/Z2mU03RKpm8c1C40PBFqto6xH4pHv4D03aXYjUBnuWS2R1i8X13l2ZOL4gsaLeB4z1mM",
	"gN4RzqhCmqAX5DwJJsGDUcFmcnT0YAxgM3nIGJebh4zDgnzZBL3gDnOC50n9aoMxK92TUvZg+3sc4hWW",
	"/YilQa/7O5bmE5Ys5/q0CRGI55QWx+sYl03ax8djLzE1skYqy+cJiQqK+ooArY7sLciXJlV9ANKcZTm6",
	"G/on0K9pj21OEDTsuqK5kjITLVLmnLeBf1NJjVaI5aYyqql5kE3r62eQ+E+l1ZU2B4l/cM19kHgJ30Fr",
	"Qs3zp8M9ldrO4Q4SltnvCn0r5enr0itzCueJtHU60JfM/1HU9ShtXZd7Z6v3fPnxvGpgIUZRrD9UHaOq",
	"S1JGe9VPbf0XYbpAlKGYcIikTiQSnajdE7mqKKJ5LlHKYkhs7CHMgILjcsJio99cb/4vAAD//6qrTLpj",
	"YAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
