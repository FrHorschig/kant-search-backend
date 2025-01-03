import argparse
import os
import glob
import html
import re


# === Helper functions ===
def find_encoding(filename):
    match = re.search(r"(\d+)\.xml$", filename)
    if match:
        n = int(match.group(1))
        if n in [10, 15, 16, 23]:
            return "iso-8859-1"
    return "utf-8"


# === Sanitation steps ===
def replace_encoding(enc, content):
    if enc == "iso-8859-1":
        content = content.replace('encoding="ISO-8859-1"', 'encoding="UTF-8"')
    return content


def remove_html_encodings(content):
    content = html.unescape(content)
    content = re.sub(r"&alpha;", "α", content)
    content = re.sub(r"&Alpha;", "Α", content)
    content = re.sub(r"&beta;", "β", content)
    content = re.sub(r"&Beta;", "Β", content)
    content = re.sub(r"&gamma;", "γ", content)
    content = re.sub(r"&Gamma;", "Γ", content)
    content = re.sub(r"&delta;", "δ", content)
    content = re.sub(r"&Delta;", "Δ", content)
    content = re.sub(r"&epsilon;", "ε", content)
    content = re.sub(r"&Epsilon;", "Ε", content)
    content = re.sub(r"&zeta;", "ζ", content)
    content = re.sub(r"&Zeta;", "Ζ", content)
    content = re.sub(r"&eta;", "η", content)
    content = re.sub(r"&Eta;", "Η", content)
    content = re.sub(r"&theta;", "θ", content)
    content = re.sub(r"&theata;", "θ", content)
    content = re.sub(r"&Theta;", "Θ", content)
    content = re.sub(r"&iota;", "ι", content)
    content = re.sub(r"&Iota;", "Ι", content)
    content = re.sub(r"&kappa;", "κ", content)
    content = re.sub(r"&Kappa;", "Κ", content)
    content = re.sub(r"&lambda;", "λ", content)
    content = re.sub(r"&Lambda;", "Λ", content)
    content = re.sub(r"&my;", "μ", content)
    content = re.sub(r"&My;", "Μ", content)
    content = re.sub(r"&ny;", "ν", content)
    content = re.sub(r"&Ny;", "Ν", content)
    content = re.sub(r"&xi;", "ξ", content)
    content = re.sub(r"&Xi;", "Ξ", content)
    content = re.sub(r"&omikron;", "ο", content)
    content = re.sub(r"&Omikron;", "Ο", content)
    content = re.sub(r"&pi;", "π", content)
    content = re.sub(r"&Pi;", "Π", content)
    content = re.sub(r"&rho;", "ρ", content)
    content = re.sub(r"&Rho;", "Ρ", content)
    content = re.sub(r"&sigma;", "σ", content)
    content = re.sub(r"&sigma2;", "ς", content)
    content = re.sub(r"&Sigma;", "Σ", content)
    content = re.sub(r"&tau;", "τ", content)
    content = re.sub(r"&Tau;", "Τ", content)
    content = re.sub(r"&ypsilon;", "υ", content)
    content = re.sub(r"&Ypsilon;", "Υ", content)
    content = re.sub(r"&phi;", "φ", content)
    content = re.sub(r"&Phi;", "Φ", content)
    content = re.sub(r"&chi;", "χ", content)
    content = re.sub(r"&Chi;", "Χ", content)
    content = re.sub(r"&psi;", "ψ", content)
    content = re.sub(r"&Psi;", "Ψ", content)
    content = re.sub(r"&omega;", "ω", content)
    content = re.sub(r"&Omega;", "Ω", content)
    return content


# === Combining everything ===
def process_file(in_file):
    enc = find_encoding(in_file)
    with open(in_file, "r", encoding=enc) as file:
        content = file.read()
    content = replace_encoding(enc, content)
    content = remove_html_encodings(content)
    return content


def write_file(path):
    with open(path, "w", encoding="utf-8") as file:
        file.write(content)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "read_path", help="The directory path containing the XML files to read."
    )
    parser.add_argument(
        "write_path", help="The directory path to write the updated XML files to."
    )
    read_path = parser.parse_args().read_path
    write_path = parser.parse_args().write_path

    files = glob.glob(os.path.join(read_path, "*.xml"))
    for file in files:
        try:
            content = process_file(file)
            file_name = os.path.basename(file)
            write_file(os.path.join(write_path, file_name))
        except Exception as e:
            print("Error in ", file, ":")
            print(e)
