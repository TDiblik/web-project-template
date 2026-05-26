import {zodResolver} from "@hookform/resolvers/zod";
import {t} from "i18next";
import * as z from "zod";

export {z, zodResolver};

// ----------- General -----------
export const EmailSchema = z.email({error: () => t("validation.email.invalid")});

export const PasswordSchema = z
  .string({error: () => t("validation.password.required")})
  .min(6, {error: () => t("validation.password.minLength")})
  .max(32, {error: () => t("validation.password.maxLength")})
  .regex(/[A-Z]/, {error: () => t("validation.password.uppercase")})
  .regex(/[a-z]/, {error: () => t("validation.password.lowercase")})
  .regex(/[0-9]/, {error: () => t("validation.password.number")})
  .refine((val) => new TextEncoder().encode(val).length <= 72, {error: () => t("validation.password.tooLong")});
export const ConfirmPasswordSchema = z.string({error: () => t("validation.confirmPassword.required")});
// -----------------------------

// ----------- Pages -----------
export const LoginFirstPageSchema = z.object({
  email: EmailSchema,
  password: PasswordSchema,
});
export type LoginPageFormType = z.infer<typeof LoginFirstPageSchema>;

export const SignUpFirstPageSchema = z
  .object({
    email: EmailSchema,
    password: PasswordSchema,
    confirmPassword: ConfirmPasswordSchema,
    useUsername: z.boolean(),
    firstName: z.string().optional(),
    lastName: z.string().optional(),
    username: z.string().optional(),
  })
  .refine(({useUsername, username}) => (useUsername ? (username?.trim().length ?? 0) > 0 : true), {
    path: ["username"],
    error: () => t("validation.username.required"),
  })
  .refine(({useUsername, username}) => (useUsername ? (username?.trim().length ?? 0) >= 3 : true), {
    path: ["username"],
    error: () => t("validation.username.minLength"),
  })
  .refine(({useUsername, username}) => (useUsername ? (username?.trim().length ?? 0) <= 50 : true), {
    path: ["username"],
    error: () => t("validation.username.maxLength"),
  })
  .refine(({useUsername, firstName}) => (!useUsername ? (firstName?.trim().length ?? 0) > 0 : true), {
    path: ["firstName"],
    error: () => t("validation.firstName.required"),
  })
  .refine(({useUsername, lastName}) => (!useUsername ? (lastName?.trim().length ?? 0) > 0 : true), {
    path: ["lastName"],
    error: () => t("validation.lastName.required"),
  })
  .refine(({password, confirmPassword}) => password === confirmPassword, {
    path: ["confirmPassword"],
    error: () => t("validation.confirmPassword.mismatch"),
  });
export type SignUpPageFormType = z.infer<typeof SignUpFirstPageSchema>;
export type LoginOrSignUpPageFormType = LoginPageFormType | SignUpPageFormType;
