package options

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.YunjingOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
