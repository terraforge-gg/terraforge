import z from "zod";

export const signInSchema = z.object({
  username: z.string(),
  password: z.string(),
});

const USERNAME_MIN_LENGTH = 3;
const PASSWORD_MIN_LENGTH = 8;
const PASSWORD_POLICY_MESSAGE = `Password must include a special character.`;

export const signUpSchema = z.object({
  username: z.string().min(USERNAME_MIN_LENGTH, {
    message: `Username must be more than ${USERNAME_MIN_LENGTH} characters.`,
  }),
  email: z.email({
    message: "Please enter a valid email address",
  }),
  password: z
    .string()
    .min(PASSWORD_MIN_LENGTH, {
      message: `Password must be at least ${PASSWORD_MIN_LENGTH} characters long`,
    })
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})/, {
      message: PASSWORD_POLICY_MESSAGE,
    }),
});
