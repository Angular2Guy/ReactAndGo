const esbuild = require("esbuild");
const { sassPlugin } = require("esbuild-sass-plugin");
/*
esbuild
  .build({
    entryPoints: ["src/Application.tsx", "src/application.scss"],
    outdir: "public/assets",
    bundle: true,    
    minify: true,    
    plugins: [sassPlugin()],
  })
  .then(() => console.log("⚡ Build complete! ⚡"))
  .catch(() => process.exit(1));
  */

async function startWatching() {
    const ctx = await esbuild.context({
        entryPoints: ["src/Application.tsx", "src/application.scss"],
        outdir: "public/assets",
        bundle: true,    
        minify: true,    
        treeShaking: true,
        plugins: [sassPlugin()],    
      });
    await ctx.watch();
}  
  
  startWatching();