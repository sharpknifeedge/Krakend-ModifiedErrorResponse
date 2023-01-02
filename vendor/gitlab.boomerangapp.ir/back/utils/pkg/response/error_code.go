package response

const call = "\n با تیم پشتیبانی تماس بگیرید."
const callOnRepeat = "\n در صورت وقوع مجدد خطا کمی صبر کنید سپس" + call

const (
	//Success on message success status: 200,201
	Success uint16 = iota
	//UnknownError other errors like incorrect username password
	UnknownError
	//NotFound http status: 404
	NotFound
	//UnprocessableEntity http status: 422
	UnprocessableEntity
	//AccessDenied http status: 403
	AccessDenied
	//Unauthorized http status: 401
	Unauthorized
	//InternalError on Database or other internal errors http status: 502,503
	InternalError
	//ConnectionError When the connection to a broker fails( bad gateway ), http status: 502
	ConnectionError
	//ServiceDisabled status:+500, 503
	ServiceDisabled
	//TooManyRequests Too Many Requests http status: 429
	TooManyRequests
)

var messages = []Message{
	{
		Header: "با موفقیت انجام شد",
	},
	{
		Header: "خطا",
		Body:   call,
	},
	{
		Header: "چیزی یافت نشد",
		Body:   call,
	},
	{
		Header: "خطا در پردازش ورودی",
		Body:   call,
	},
	{
		Header: "دسترسی ندارید",
		Body:   "شما به این بخش دسترسی ندارید" + call,
	},
	{
		Header: "هویت شما شناسایی نشد",
		Body:   "لطفا مجددا وارد شوید" + call,
	},
	{
		Header: "خطای داخلی",
		Body:   "خطا در ارتباط با پایگاه داده یا سایر خطا های داخلی" + callOnRepeat,
	},
	{
		Header: "ناتوان در اتصال",
		Body:   "سیستم موقتا در اتصال به کارگزار مربوطه با خطا مواجه شده است" + callOnRepeat,
	},
	{
		Header: "سرویس موقتا غیر فعال است",
		Body:   "سرویس درخواستی شما موقتا غیر فعال است.",
	},
	{
		Header: "محدودیت زمانی",
		Body:   "دسترسی شما به علت تکرار عملیات بیش از حد مجاز موقتا غیر فعال شده است. چند دقیقه صبر کنید،" + callOnRepeat,
	},
}
