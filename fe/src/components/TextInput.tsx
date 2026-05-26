import {motion} from "motion/react";
import type {AnimationEventHandler} from "react";
import {Controller, useFormContext} from "react-hook-form";
import type {FormFieldProps} from "../utils/form";

export type TextInputProps = FormFieldProps & {
  placeholder?: string;
  containerProps?: React.HTMLAttributes<HTMLDivElement>;
  inputProps?: React.InputHTMLAttributes<HTMLInputElement>;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
  labelSpanProps?: React.HTMLAttributes<HTMLSpanElement>;
  labelSpanAdditionalClassname?: string;
  errorSpanProps?: AnimationEventHandler<HTMLSpanElement>;
  errorSpanPropsAdditionalClassname?: string;
  isOptional?: boolean;
  hasBigText?: boolean;
};

export const TextInput: React.FC<TextInputProps> = (props) => {
  const form = useFormContext();

  return (
    <Controller
      name={props.name}
      control={form.control}
      render={({field, fieldState}) => {
        const {ref, ...rest} = field;
        const hasError = !!fieldState.error;

        return (
          <div className="form-control w-full" ref={ref} {...props.containerProps}>
            <fieldset className="fieldset">
              <legend className="fieldset-legend pb-1">
                {props.label && (
                  <label htmlFor={props.name} className={`${props.hasBigText && "text-base font-normal"} label`} {...props.labelProps}>
                    <span className={`label-text ${hasError ? "text-red-500" : ""} ${props.labelSpanAdditionalClassname}`} {...props.labelSpanProps}>
                      {props.label}
                    </span>
                  </label>
                )}
              </legend>
              <input
                id={props.name}
                type="text"
                placeholder={props.placeholder}
                disabled={props.isDisabled}
                className={`input input-bordered w-full ${hasError ? "input-error" : ""}`}
                {...props.inputProps}
                {...rest}
              />
              {props.isOptional && <p className="label">Optional</p>}
            </fieldset>

            {hasError && (
              <motion.span
                initial={{opacity: 0}}
                animate={{opacity: 1}}
                className={`text-red-500 text-sm text-center mb-2 ${props.errorSpanPropsAdditionalClassname}`}
                {...props.errorSpanProps}
              >
                {fieldState.error?.message}
              </motion.span>
            )}
          </div>
        );
      }}
    />
  );
};
