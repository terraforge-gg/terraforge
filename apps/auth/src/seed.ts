import { auth } from "./auth.js";
import { env } from "./env.js";

const main = async () => {
  try {
    const response = await auth.api.isUsernameAvailable({
      body: {
        username: env.SEED_USER_USERNAME,
      },
    });

    if (!response?.available) {
      console.log(`User '${env.SEED_USER_USERNAME}' already exists.`);
      return;
    }

    const res = await auth.api.signUpEmail({
      body: {
        name: env.SEED_USER_USERNAME,
        username: env.SEED_USER_USERNAME,
        email: env.SEED_USER_EMAIL,
        password: env.SEED_USER_PASSWORD,
        displayUsername: env.SEED_USER_USERNAME,
      },
    });

    console.log(res.user.id);
  } catch (error) {
    console.log("Failed to seed user. Error: ", error);
  }
};

await main();
process.exit(0);
