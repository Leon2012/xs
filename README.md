"#xs" 

###字节序###

xs网络协议采用主机字节序(小端字节序)


###php和golang封包方法###
```
封包:
php:
pack('CCCCI', $this->cmd, $this->arg1, $this->arg2, strlen($this->buf1), strlen($this->buf)) . $this->buf . $this->buf1;
//pack("CCCCI") = uchar(1) + uchar(1) + uchar(1) + uchar(1) + uint(4) = 8

golang:
bufBytes := []byte(x.Buf)
buf1Bytes := []byte(x.Buf1)
bufLen := 8 + len(bufBytes) + len(buf1Bytes)
var buf []byte
buf = make([]byte, bufLen)
var cmd uint8
var arg1 uint8
var arg2 uint8
cmd = uint8(x.Cmd)
arg1 = uint8(x.Arg1)
arg2 = uint8(x.Arg2)
buf[0] = byte(cmd)
buf[1] = byte(arg1)
buf[2] = byte(arg2)

buf1Len := len(x.Buf1)
buf1LenByte := byte(buf1Len)
buf[3] = buf1LenByte

bufLen = len(x.Buf)
bufLenBytes := make([]byte, 4)
//copy(buf[4:8], Uint32ToBytes(uint32(bufLen)))
PutUInt32(uint32(bufLen), bufLenBytes, 60)
copy(buf[4:8], bufLenBytes)

copy(buf[8:(8+len(bufBytes))], bufBytes)

idx := 8 + len(bufBytes)
copy(buf[idx:(idx+len(buf1Bytes))], buf1Bytes)

解包:
php:
$resFormat = 'Idocid/Irank/Iccount/ipercent/fweight';
$metas = unpack(self::$_resFormat, $p);

golang:
metas := make(map[string]interface{})
docIdBytes := make([]byte, 4)
copy(docIdBytes, data[0:4])
fmt.Println(docIdBytes)
docId := binary.LittleEndian.Uint32(docIdBytes)
metas["docid"] = docId

rankBytes := make([]byte, 4)
copy(rankBytes, data[4:8])
rank := binary.LittleEndian.Uint32(rankBytes)
metas["rank"] = rank

ccountBytes := make([]byte, 4)
copy(ccountBytes, data[8:12])
ccount := binary.LittleEndian.Uint32(ccountBytes)
metas["ccount"] = ccount

percentBytes := make([]byte, 4)
copy(percentBytes, data[12:16])
percent := int32(binary.LittleEndian.Uint32(percentBytes))
metas["percent"] = percent

weightBytes := make([]byte, 4)
copy(weightBytes, data[16:20])
weight := ByteToFloat32(weightBytes)
metas["weight"] = weight

I对应的是uint32类型占4个字节，对应的golang中就是[4]byte。
i对应的是int32
C对应的是uchar类型占1个字节，对应golang中就是[1]byte
```

###协议###
