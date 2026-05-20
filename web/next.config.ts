import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    reactStrictMode: true,
    reactCompiler: true,
    cacheComponents: false,
    output: "standalone",
};

export default nextConfig;
