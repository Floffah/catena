import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    reactStrictMode: true,
    reactCompiler: true,
    cacheComponents: false,

    rewrites: async () => [
        {
            source: "/api/:path*",
            destination: `${process.env.NEXT_PUBLIC_CATENA_INSTANCE_URL}/:path*`,
        },
    ],
};

export default nextConfig;
