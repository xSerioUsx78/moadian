# How to use it

Opening the private key file.

```
privateKey, err := os.ReadFile("private.prv")
if err != nil {
    panic(err)
}
```

Then create an instance of the Moadian interface and pass the information.

```
var moadian service.Moadian
baseInfo := service.MoadiInfoStruct{
    FiscalID: "Your Fiscal id",
    PrivateKey: privateKey,
    TaxServerPublicKey: "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAxdzREOEfk3vBQogDPGTMqdDQ7t0oDhuKMZkA+Wm1lhzjjhAGfSUOuDvOKRoUEQwP8oUcXRmYzcvCUgcfoRT5iz7HbovqH+bIeJwT4rmLmFcbfPke+E3DLUxOtIZifEXrKXWgSVPkRnhMgym6UiAtnzwA1rmKstJoWpk9Nv34CYgTk8DKQN5jQJqb9L/Ng0zOEEtI3zA424tsd9zv/kP4/SaSnbbnj0evqsZ29X6aBypvnTnwH9t3gbWM4I9eAVQhPYClawHTqvdaz/O/feqfm06QBFnCgL+CBdjLs30xQSLsPICjnlV1jMzoTZnAabWP6FRzzj6C2sxw9a/WwlXrKn3gldZ7Ctv6Jso72cEeCeUI1tzHMDJPU3Qy12RQzaXujpMhCz1DVa47RvqiumpTNyK9HfFIdhgoupFkxT14XLDl65S55MF6HuQvo/RHSbBJ93FQ+2/x/Q2MNGB3BXOjNwM2pj3ojbDv3pj9CHzvaYQUYM1yOcFmIJqJ72uvVf9Jx9iTObaNNF6pl52ADmh85GTAH1hz+4pR/E9IAXUIl/YiUneYu0G4tiDY4ZXykYNknNfhSgxmn/gPHT+7kL31nyxgjiEEhK0B0vagWvdRCNJSNGWpLtlq4FlCWTAnPI5ctiFgq925e+sySjNaORCoHraBXNEwyiHT2hu5ZipIW2cCAwEAAQ==",
    TaxServerKeyID: "6a2bcd88-a871-4245-a393-2843eafe6e02",
}
moadian = &baseInfo
```

Get the token and pass the token and a map of invoice data to SendInvoice method.

```
tokenResData, err := moadian.GetToken()
if err != nil {
    panic(err)
}
token := tokenResData.Result.Data.Token
packetData := map[string]any{
    "header": map[string]any{
        "taxid": moadian.GenerateTaxID(time.Now(), 1),
        "indatim": 1692085800000,
        "indati2m": 1692085803000,
        "inty": 1,
        "inno": "0000001140",
        "irtaxid": nil,
        "inp": 1,
        "ins": 1,
        "tins" : "10101704295",
        "tob" : 2,
        "bid" : nil,
        "tinb" : nil,
        "sbc" : nil,
        "bpc" : nil,
        "bbc" : nil,
        "ft" : nil,
        "bpn" : nil,
        "scln" : nil,
        "scc" : nil,
        "crn" : nil,
        "billid" : nil,
        "tprdis" : 61000000,
        "tdis" : 0,
        "tadis" : 61000000,
        "tvam" : 5490000,
        "todam" : 0,
        "tbill" : 66490000,
        "tonw": nil,
        "torv": nil,
        "tocv": nil,
        "setm" : 2,
        "cap" : nil,
        "insp" : 61000000,
        "tvop" : 5490000,
        "tax17" : nil,
    },
    "body": []map[string]any{
        {
            "sstid" : nil,
            "sstt": "FooBar",
            "am" : 5,
            "mu": nil,
            "fee" : 12200000,
            "cfee" : nil,
            "cut" : nil,
            "exr" : nil,
            "ssrv": nil,
            "sscv": nil,
            "prdis" : 61000000,
            "dis" : 0,
            "adis" : 61000000,
            "vra" : 9,
            "vam" : 5490000,
            "odt" : nil,
            "odr" : nil,
            "odam" : nil,
            "olt" : nil,
            "olr" : nil,
            "olam" : nil,
            "consfee" : nil,
            "spro" : nil,
            "bros" : nil,
            "tcpbs" : nil,
            "cop" : nil,
            "vop" : 5490000,
            "bsrn" : nil,
            "tsstam" : 66490000,
        },
    },
    "payments": []map[string]any{},
}

sendInvoiceResData, err := moadian.SendInvoice(
    &token,
    &packetData,
)
refNumber := sendInvoiceResData.Result[0].ReferenceNumber
```

And you can also inquiry the invoice by ReferenceNumber

```
refNumbers := []string{refNumber}
inquiryResData, err := moadian.InquiryByReferenceNumber(&token, &refNumbers)
fmt.Println(inquiryResData.Result.Data[0].Status)
fmt.Println(inquiryResData.Result.Data[0].Data.Error)
fmt.Println(inquiryResData.Result.Data[0].Data.Warning)
```
