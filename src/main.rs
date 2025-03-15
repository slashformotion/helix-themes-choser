use askama::Template; // bring trait in scope
use crossbeam::scope;
use once_cell::sync::Lazy;
use std::collections::HashMap;
use std::io::Write;
use std::process::Command;
use tempfile;

const LANGUAGE_FOLDER: &str = "./languages/";
const OUT_FOLDER: &str = "./out/";
static LANGUAGES: Lazy<HashMap<&'static str, &'static str>> = Lazy::new(|| {
    let mut map = HashMap::new();
    map.insert("golang", "main.go");
    map.insert("python3", "main.py");
    map
});

#[derive(Template)] // this will generate the code...
#[template(path = "recipe.tape", escape = "none")] // using the template in this path, relative
struct TapeFile<'a> {
    // the name of the struct can be anything
    out: &'a str,   // the field name should match the variable name
    theme: &'a str, // the field name should match the variable name
    file: &'a str,  // the field name should match the variable name
                    // in your template
}

const THEMES: [&str; 5] = [
    "rasmus",
    "vintage",
    "starlight",
    "gruber-darker",
    "solarized_dark",
    // "github_light",
    // "ayu_evolve",
    // "material_deep_ocean",
    // "gruvbox_dark_soft",
    // "catppuccin_frappe",
    // "modus_vivendi",
    // "rose_pine_moon",
    // "github_light_colorblind",
    // "everblush",
    // "naysayer",
    // "hex_toxic",
    // "zenburn",
    // "eiffel",
    // "tokyonight_moon",
    // "catppuccin_macchiato",
    // "hex_steel",
    // "nord-night",
    // "github_dark_colorblind",
    // "base16_terminal",
    // "monokai_soda",
    // "seoul256-dark",
    // "seoul256-light-hard",
    // "term16_light",
    // "tokyonight",
    // "seoul256-dark-hard",
    // "iceberg-light",
    // "ferra",
    // "cyan_light",
    // "bogster",
    // "yo_light",
    // "base16_default_light",
    // "bogster_light",
    // "beans",
    // "iroaseta",
    // "gruvbox_light_hard",
    // "tokyonight_day",
    // "monokai_aqua",
    // "monokai_pro_machine",
    // "horizon-dark",
    // "molokai",
    // "monokai_pro",
    // "snazzy",
    // "kaolin-valley-dark",
    // "acme",
    // "flexoki_light",
    // "material_darker",
    // "hex_poison",
    // "darcula-solid",
    // "autumn",
    // "github_dark_tritanopia",
    // "ingrid",
    // "ayu_light",
    // "github_light_tritanopia",
    // "zed_onelight",
    // "catppuccin_latte",
    // "zed_onedark",
    // "tokyonight_storm",
    // "noctis",
    // "yellowed",
    // "base16_default_dark",
    // "monokai_pro_octagon",
    // "hex_lavender",
    // "jetbrains_dark",
    // "ttox",
    // "github_dark_dimmed",
    // "varua",
    // "rose_pine_dawn",
    // "modus_vivendi_tritanopia",
    // "papercolor-light",
    // "material_palenight",
    // "kaolin-dark",
    // "vim_dark_high_contrast",
    // "seoul256-light",
    // "heisenberg",
    // "fleet_dark",
    // "voxed",
    // "seoul256-light-soft",
    // "nightfox",
    // "noctis_bordo",
    // "merionette",
    // "nord_light",
    // "seoul256-dark-soft",
    // "solarized_light",
    // "serika-light",
    // "gruvbox_light_soft",
    // "modus_operandi_deuteranopia",
    // "github_dark_high_contrast",
    // "modus_vivendi_tinted",
    // "catppuccin_mocha",
    // "yo",
    // "ayu_mirage",
    // "monokai_pro_ristretto",
    // "autumn_night",
    // "everforest_light",
    // "material_oceanic",
    // "flexoki_dark",
    // "mellow",
    // "meliora",
    // "onedarker",
    // "flatwhite",
    // "kanagawa",
    // "gruvbox_dark_hard",
    // "github_light_high_contrast",
    // "modus_operandi_tritanopia",
    // "adwaita-dark",
    // "poimandres_storm",
    // "darcula",
    // "iceberg-dark",
    // "jellybeans",
    // "amberwood",
    // "everforest_dark",
    // "modus_vivendi_deuteranopia",
    // "serika-dark",
    // "curzon",
    // "onedark",
    // "carbonfox",
    // "poimandres",
    // "dark_plus",
    // "onelight",
    // "boo_berry",
    // "ayu_dark",
    // "dark_high_contrast",
    // "pop-dark",
    // "monokai_pro_spectrum",
    // "emacs",
    // "dracula",
    // "kaolin-light",
    // "modus_operandi",
    // "term16_dark",
    // "gruvbox",
    // "new_moon",
    // "adwaita-light",
    // "dracula_at_night",
    // "github_dark",
    // "ao",
    // "kanagawa-dragon",
    // "spacebones_light",
    // "papercolor-dark",
    // "doom_acario_dark",
    // "night_owl",
    // "modus_operandi_tinted",
    // "base16_transparent",
    // "monokai",
    // "nord",
    // "sonokai",
    // "yo_berry",
    // "gruvbox_light",
    // "rose_pine",
    // "penumbra+",
    // "sunset",
];

fn create_screenshot(lang: &str, language_file: &str, theme: &str) {
    let out = format!("{}{}-{}", OUT_FOLDER, theme, language_file);
    let file = format!("{}{}", LANGUAGE_FOLDER, &language_file);
    let t = TapeFile {
        out: &out,
        file: &file,
        theme: &theme,
    };
    let mut temp = tempfile::NamedTempFile::new().unwrap();
    let rend = t.render().unwrap();
    temp.write_all(rend.as_bytes()).unwrap();

    match Command::new("sh")
        .arg("-c")
        .arg(format!("vhs {}", &temp.path().to_string_lossy()))
        .output()
    {
        Ok(_) => println!("success: lang={}, theme={}", lang, theme),
        Err(e) => eprintln!("success: lang={}, theme={}: ERROR={:?}", lang, theme, e),
    }
}

fn main() {
    scope(|s| {
        for theme in THEMES.iter() {
            for (lang, language_file) in LANGUAGES.iter() {
                s.spawn(move |_| {
                    create_screenshot(lang, language_file, theme);
                });
            }
        }
    })
    .unwrap();
}
