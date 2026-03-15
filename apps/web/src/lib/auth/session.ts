import { cookies } from "next/headers";
import { env } from "@/env";

export type SessionUser = {
  id: string;
  name: string;
  email: string;
  emailVerified: boolean;
  image?: string | null;
  createdAt: string;
  updatedAt: string;
  username: string;
  displayUsername: string;
};

export type Session = {
  id: string;
  expiresAt: string;
  token: string;
  createdAt: string;
  ipAddress: string;
  userAgent: string;
  userId: string;
};

export type SessionResponse = {
  session: Session;
  user: SessionUser;
};

const getServerSession = async (): Promise<SessionResponse | null> => {
  const cookieStore = await cookies();
  const res = await fetch(env.NEXT_PUBLIC_AUTH_URL + "/api/auth/get-session", {
    method: "GET",
    headers: {
      Cookie: cookieStore.toString(),
    },
  });

  if (!res.ok) {
    return null;
  }

  const body = await res.text();
  if (body == "null") {
    return null;
  }

  return JSON.parse(body) as SessionResponse;
};

export default getServerSession;
