package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"reflect"

	"github.com/pkg/errors"

	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}
func GetMsgID(sendID string) string {
	t := Int64ToString(GetCurrentTimestampByNano())
	return Md5(t + sendID + Int64ToString(rand.Int63n(GetCurrentTimestampByNano())))
}
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}

//Get the current timestamp by Second

func GetCurrentTimestampBySecond() int64 {
	return time.Now().Unix()
}

// GetCurrentTimestampByMill Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

// UnixNanoSecondToTime Convert nano timestamp to time.Time type
func UnixNanoSecondToTime(nanoSecond int64) time.Time {
	return time.Unix(0, nanoSecond)
}

// GetCurrentTimestampByNano Get the current timestamp by Nano
// GetCurrentTimestampByNano 返回当前时间戳，单位为纳秒。
// 该函数使用 time.Now().UnixNano() 方法获取当前的 Unix 时间戳，
// 并将其以纳秒为单位返回。这对于需要高精度时间戳的场景非常有用。
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

// StructToJsonString 将结构体转换为JSON字符串。
// 该函数接受一个接口类型的参数，允许它接受任何类型的结构体。
// 它的主要作用是将传入的结构体实例序列化为JSON格式的字符串。
// 参数:
//
//	param 接口类型，代表任意类型的结构体。
//
// 返回值:
//
//	string 序列化后的JSON字符串。
//	如果传入的参数不是结构体，或者JSON序列化失败，函数将panic。
func StructToJsonString(param interface{}) string {
	// 使用json.Marshal尝试将param转换为JSON格式的字节切片。
	// 这一步可能会出错，比如如果param不是结构体，或者包含无法序列化的字段。
	dataType, err := json.Marshal(param)
	if err != nil {
		// 如果发生错误，函数将panic，终止执行并抛出运行时错误。
		panic(err)
	}
	// 将序列化后的JSON字节切片转换为字符串格式。
	dataString := string(dataType)
	// 返回转换后的JSON字符串。
	return dataString
}

func StructToJsonStringDefault(param interface{}) string {
	if reflect.TypeOf(param).Kind() == reflect.Slice && reflect.ValueOf(param).Len() == 0 {
		return "[]"
	}
	return StructToJsonString(param)
}

// JsonStringToStruct The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	return Wrap(json.Unmarshal([]byte(s), args), "json Unmarshal failed")
}
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

//Convert timestamp to time.Time type

// UnixSecondToTime 将Unix时间戳（秒）转换为time.Time类型。
// 该函数接收一个int64类型的参数second，表示从1970年1月1日零时（Unix纪元的开始）到指定时间的秒数。
// 返回值是一个time.Time类型的值，表示转换后的日期和时间。
// 注意：此函数不处理纳秒部分，因此将纳秒部分设置为0。
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// IntToString 将整数转换为字符串。
// 这个函数使用strconv.FormatInt方法将整数i转换为字符串格式。
// 参数i是需要转换的整数。
// 返回值是转换后的字符串。
func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// Int32ToString 将32位有符号整数转换为字符串。
// 这个函数使用strconv标准库中的FormatInt函数来实现转换。
// 参数:
//
//	i: 需要转换的32位有符号整数。
//
// 返回值:
//
//	转换后的字符串。
func Int32ToString(i int32) string {
	// 使用strconv.FormatInt函数将int32类型的i转换为字符串
	return strconv.FormatInt(int64(i), 10)
}

// Int64ToString 将 int64 类型的数字转换为字符串。
// 该函数使用 strconv 标准库中的 FormatInt 函数进行转换。
// 参数:
//
//	i: 需要转换的 int64 类型的数字。
//
// 返回值:
//
//	转换后的字符串。
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt64 将字符串转换为int64类型。
// 参数:
//
//	i: 待转换的字符串。
//
// 返回值:
//
//	转换后的int64类型值。
//
// 说明:
//
//	本函数使用strconv.ParseInt进行转换，忽略了错误处理，因为预期输入总是可转换的。
func StringToInt64(i string) int64 {
	// ParseInt按照指定的基数和位数将字符串转换为int64。
	// 这里忽略了转换过程中的可能错误，因为应该保证输入的字符串总是可以转换的。
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}

// StringToInt 将字符串转换为整数。
// 这个函数接收一个字符串作为参数，尝试将其转换为整数类型并返回。
// 如果转换失败，将返回一个默认的整数值（在本例中未明确处理错误情况）。
//
// 参数:
//
//	i: 待转换的字符串。
//
// 返回值:
//
//	转换后的整数。如果转换失败，则返回一个默认值（本例中未处理）。
func StringToInt(i string) int {
	// 使用strconv.Atoi函数尝试将字符串转换为整数。
	// Atoi函数返回两个值：转换后的整数和一个错误值。
	// 在这里，我们只关心转换后的整数，因此错误值被忽略。
	j, _ := strconv.Atoi(i)

	// 返回转换后的整数。
	// 注意：如果转换失败，'j' 的值将是0，这种情况下调用者需要知道0是一个有效值还是转换错误的结果。
	return j
}

// RunFuncName 返回调用该函数的上级函数的名称。
// 该函数通过runtime包提供的Caller函数获取调用栈信息，然后通过FuncForPC函数获取被调用函数的名称。
// 使用runtime包处理是因为Go没有直接内置获取当前函数名的语言特性。
// 需要注意的是，本函数只返回上级函数的名称，不包括参数等其他信息。
func RunFuncName() string {
	// 获取调用者的调用信息，这里传递的参数2表示获取调用RunFuncName函数的上级函数的信息。
	// runtime.Caller函数返回四个值，但我们只关心第一个，即程序计数器（program counter, PC）。
	pc, _, _, _ := runtime.Caller(2)

	// 根据获取到的程序计数器，使用runtime.FuncForPC函数获取对应的函数信息。
	// 然后调用CleanUpfuncName函数清理函数名称，使其更易于阅读和理解。
	return CleanUpfuncName(runtime.FuncForPC(pc).Name())
}

type LogInfo struct {
	Info string `json:"info"`
}

// IsContain judge a string whether in the  string list
// IsContain 检查目标字符串是否存在于字符串列表中。
//
// target: 需要查找的目标字符串。
// List: 用于查找的字符串列表。
//
// 返回值: 如果目标字符串存在于列表中，则返回 true；否则返回 false。
func IsContain(target string, List []string) bool {
	// 遍历字符串列表
	for _, element := range List {
		// 比较当前列表项与目标字符串是否相等
		if target == element {
			// 如果相等，则目标字符串存在于列表中，返回 true
			return true
		}
	}
	// 如果遍历结束后仍未找到目标字符串，返回 false
	return false
}

// IsContainInt 检查目标整数是否在给定的整数列表中。
// 该函数通过遍历列表中的每个元素进行比较，避免了复杂的逻辑。
// 参数:
//
//	target: 需要查找的目标整数。
//	List: 包含整数的列表。
//
// 返回值:
//
//	如果目标整数在列表中，则返回true；否则返回false。
func IsContainInt(target int, List []int) bool {

	// 遍历列表中的每个元素。
	for _, element := range List {

		// 如果目标整数与当前元素相等，则返回true。
		if target == element {
			return true
		}
	}
	// 如果遍历完所有元素都没有找到相等的目标整数，则返回false。
	return false
}

func IsContainUInt32(target uint32, List []uint32) bool {

	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false

}
func GetSwitchFromOptions(Options map[string]bool, key string) (result bool) {
	if flag, ok := Options[key]; !ok || flag {
		return true
	}
	return false
}
func SetSwitchFromOptions(Options map[string]bool, key string, value bool) {
	Options[key] = value
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

// Unwrap 函数用于递归地解包错误，直到解包出最根本的错误。
// 这个函数旨在处理那些实现了 Unwrap 方法的错误类型，通过不断地调用 Unwrap 方法，
// 可以获取错误链中的下一个错误，直到没有更多的错误可以解包为止。
func Unwrap(err error) error {
	// 循环检查错误是否为 nil，如果不为 nil，则尝试解包。
	for err != nil {
		// 尝试将错误转换为具有 Unwrap 方法的接口类型。
		unwrap, ok := err.(interface {
			Unwrap() error
		})
		// 如果转换不成功，说明当前错误没有实现 Unwrap 方法，循环结束。
		if !ok {
			break
		}
		// 调用 Unwrap 方法获取下一个错误，并将其赋值给 err 变量。
		err = unwrap.Unwrap()
	}
	// 返回最终解包出的最根本的错误。
	return err
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return CleanUpfuncName(runtime.FuncForPC(pc).Name())
}
func CleanUpfuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

// StructToMap 将一个结构体对象转换为Map对象。
// 该函数接受一个interface{}类型的参数user，可以是任何实现了编码（encoding）的结构体类型。
// 参数:
//
//	user - 待转换的结构体对象，必须实现encoding。
//
// 返回值:
//
//	返回一个map[string]interface{}，其中包含了user对象的键值对表示。
//	如果转换过程中发生错误，返回nil。
func StructToMap(user interface{}) map[string]interface{} {
	// 将user对象转换为JSON字节切片。
	data, _ := json.Marshal(user)

	// 初始化一个空的map，用于存储转换后的键值对。
	m := make(map[string]interface{})

	// 将JSON字节切片转换为map。
	err := json.Unmarshal(data, &m)
	if err != nil {
		// 如果转换过程中发生错误，返回nil。
		return nil
	}

	// 返回转换后的map。
	return m
}

// KMP KMP实现字符串的模式匹配，使用KMP算法提高匹配效率。
// 参数rMainString是主字符串，rSubString是待查找的子字符串。
// 返回值isInMainString表示rSubString是否存在于rMainString中。
func KMP(rMainString string, rSubString string) (isInMainString bool) {
	// 将主字符串和子字符串转换为小写，以进行不区分大小写的匹配。
	mainString := strings.ToLower(rMainString)
	subString := strings.ToLower(rSubString)

	// 初始化主字符串和子字符串的索引。
	mainIdx := 0
	subIdx := 0

	// 获取主字符串和子字符串的长度。
	mainLen := len(mainString)
	subLen := len(subString)

	// 计算子字符串的next数组，用于KMP算法中的回溯。
	next := computeNextArray(subString)

	// 主循环进行KMP匹配。
	for {
		// 如果主字符串或子字符串的索引超出其长度，则结束循环。
		if mainIdx >= mainLen || subIdx >= subLen {
			break
		}

		// 如果当前字符匹配，则同时增加主字符串和子字符串的索引。
		if mainString[mainIdx] == subString[subIdx] {
			mainIdx++
			subIdx++
		} else {
			// 如果子字符串索引不为0，则根据next数组回溯子字符串索引。
			// 否则，只增加主字符串索引。
			if subIdx != 0 {
				subIdx = next[subIdx-1]
			} else {
				mainIdx++
			}
		}
	}

	// 如果子字符串索引等于子字符串长度，则表示找到了匹配，
	// 同时检查主字符串索引以确保没有越界。
	if subIdx >= subLen {
		if mainIdx-subLen >= 0 {
			return true
		}
	}
	return false
}

// computeNextArray
// computeNextArray 计算并返回一个用于字符串匹配的 next 数组。
// subString: 需要进行匹配的子字符串。
// 返回值：next 数组，用于指导 KMP 算法中的跳跃操作。
func computeNextArray(subString string) []int {
	// 初始化 next 数组，长度与子字符串相同。
	next := make([]int, len(subString))
	// index 用于表示子字符串中的当前位置。
	index := 0
	// i 用于遍历子字符串中的字符。
	i := 1
	// 当 i 小于子字符串的长度时，循环继续。
	for i < len(subString) {
		// 如果当前字符与 index 位置的字符相同，则更新 next 数组。
		if subString[i] == subString[index] {
			next[i] = index + 1
			i++
			index++
		} else {
			// 如果 index 不为 0，则回退 index 到 next[index-1]。
			if index != 0 {
				index = next[index-1]
			} else {
				// 如果 index 为 0，则只增加 i。
				i++
			}
		}
	}
	// 返回计算好的 next 数组。
	return next
}

// TrimStringList 去除两头的空格
func TrimStringList(list []string) (result []string) {
	for _, v := range list {
		if len(strings.Trim(v, " ")) != 0 {
			result = append(result, v)
		}
	}
	return result

}

// Intersect Get the intersection of two slices
// Intersect 计算两个切片的交集。
// 它接受两个 int64 类型的切片作为输入，并返回一个新的切片，该切片包含同时存在于两个输入切片中的所有元素。
// 参数:
//
//	slice1 - 第一个切片
//	slice2 - 第二个切片
//
// 返回值:
//
//	一个切片，包含同时存在于 slice1 和 slice2 中的所有元素。
func Intersect(slice1, slice2 []int64) []int64 {
	// 使用 map 来存储第一个切片中的元素，以便快速查找。
	m := make(map[int64]bool)
	// 初始化一个切片来存储交集结果。
	n := make([]int64, 0)

	// 遍历第一个切片，将所有元素存入 map 中。
	for _, v := range slice1 {
		m[v] = true
	}

	// 遍历第二个切片，检查元素是否在 map 中存在。
	// 如果存在，则说明该元素同时存在于两个切片中，将其添加到结果切片中。
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}

	// 返回交集结果。
	return n
}

// DifferenceSubset Get the diff of two slices
// DifferenceSubset 计算两个切片的差集。
// 它接受两个 []int64 类型的切片：mainSlice 和 subSlice，
// 并返回一个切片，其中包含所有在 mainSlice 中但不在 subSlice 中的元素。
func DifferenceSubset(mainSlice, subSlice []int64) []int64 {
	// 使用 map 来存储 subSlice 中的元素，以便快速检查元素是否存在。
	m := make(map[int64]bool)
	// n 用于存储差集的结果。
	n := make([]int64, 0)

	// 遍历 subSlice，将元素存在标记在 map 中。
	for _, v := range subSlice {
		m[v] = true
	}

	// 遍历 mainSlice，如果元素不在 subSlice 中（即 map 中没有该元素），
	// 则将其添加到结果切片 n 中。
	for _, v := range mainSlice {
		if !m[v] {
			n = append(n, v)
		}
	}

	// 返回差集结果。
	return n
}

// DifferenceSubsetString 从 mainSlice 中移除 subSlice 中包含的字符串，返回两者的差集
func DifferenceSubsetString(mainSlice, subSlice []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, v := range subSlice {
		m[v] = true
	}
	for _, v := range mainSlice {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}
func JsonDataOne(pb proto.Message) map[string]interface{} {
	return ProtoToMap(pb, false)
}

// ProtoToMap 将 Protocol Buffers 消息转换为 map[string]interface{}。
// 此函数接收一个 Protocol Buffers 消息实例和一个布尔标志 idFix。
// 如果 idFix 为 true，并且消息中包含 "id" 字段，则会将 "id" 字段的值移动到 "_id" 字段，并删除 "id" 字段。
// 参数:
//   - pb: 要转换的 Protocol Buffers 消息实例。
//   - idFix: 是否进行 ID 字段的修复，默认为 false。
//
// 返回值:
//   - 一个包含 Protocol Buffers 消息数据的 map[string]interface{}。
//   - 如果在转换过程中发生错误（例如 JSON 解析失败），则返回 nil。
func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	// 创建 jsonpb.Marshaler 实例，用于将 Protocol Buffers 消息转换为 JSON 字符串。
	// 配置 OrigName 为 true 以保留字段的原始名称。
	// 配置 EnumsAsInts 为 false 以将枚举作为字符串处理。
	// 配置 EmitDefaults 为 true 以确保所有字段都包含在输出中，即使它们是默认值。
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}

	// 将 Protocol Buffers 消息转换为 JSON 字符串。
	s, _ := marshaler.MarshalToString(pb)

	// 初始化一个空的 map 用于存储转换后的 JSON 数据。
	out := make(map[string]interface{})

	// 将 JSON 字符串解析到 out map 中。
	err := json.Unmarshal([]byte(s), &out)
	if err != nil {
		// 如果解析过程中发生错误，则返回 nil。
		return nil
	}

	// 如果 idFix 为 true，尝试将 "id" 字段的值移动到 "_id" 字段，并删除 "id" 字段。
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}

	// 返回转换后的 map。
	return out
}

func GetUserIDForMinSeq(userID string) string {
	return "u_" + userID
}

func GetGroupIDForMinSeq(groupID string) string {
	return "g_" + groupID
}

// TimeStringToTime 将字符串形式的时间转换为time.Time类型。
// 参数timeString是表示时间的字符串，格式为"2006-01-02"。
// 返回值是转换后的时间对象和一个可能的错误。
// 如果时间字符串不符合预期格式，将返回错误。
func TimeStringToTime(timeString string) (time.Time, error) {
	// time.Parse根据给定的格式和时间字符串解析时间，返回time.Time类型的时间对象和一个可能的错误。
	t, err := time.Parse("2006-01-02", timeString)
	return t, err
}

// TimeToString 将时间对象转换为字符串。
// 该函数接受一个 time.Time 类型的参数，并将其转换为 "2006-01-02" 格式的日期字符串。
// 参数:
//
//	t: 需要转换的时间对象。
//
// 返回值:
//
//	返回转换后的日期字符串。
func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

// Uint32ListConvert 将一个无符号32位整型切片转换为64位整型切片。
// 这个函数遍历输入的无符号32位整型切片，并将每个元素转换为64位整型，然后将转换后的值添加到新的切片中。
// 参数:
//
//	list - 一个无符号32位整型切片，包含需要转换的整数。
//
// 返回值:
//
//	一个64位整型切片，其中包含了从输入的32位整型切片转换而来的整数。
func Uint32ListConvert(list []uint32) []int64 {
	// 初始化一个新的64位整型切片，用于存放转换后的整数。
	var result []int64
	// 遍历输入的32位整型切片。
	for _, v := range list {
		// 将32位无符号整数转换为64位有符号整数，并添加到结果切片中。
		result = append(result, int64(v))
	}
	// 返回转换后的64位整型切片。
	return result
}
