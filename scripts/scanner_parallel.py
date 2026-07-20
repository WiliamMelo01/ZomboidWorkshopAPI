#!/usr/bin/env python3
"""
Project Zomboid Steam Workshop Parallel Scanner
Downloads Steam Workshop mods via SteamCMD and extracts ModID -> WorkshopID mappings from mod.info.
Sends extracted mappings in bulk to ZomboidWorkshopAPI.
"""

import os
import re
import sys
import json
import urllib.request
from concurrent.futures import ThreadPoolExecutor

API_URL = os.environ.get("API_URL", "http://localhost:8080/mappings/bulk")
MAX_WORKERS = int(os.environ.get("MAX_WORKERS", "4"))

def parse_mod_info(file_path, workshop_id):
    mappings = []
    if not os.path.exists(file_path):
        return mappings

    try:
        with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
            for line in f:
                line = line.strip()
                if line.startswith("id="):
                    mod_id = line.split("=", 1)[1].strip()
                    if mod_id:
                        mappings.append({"modId": mod_id, "workshopId": str(workshop_id)})
    except Exception as e:
        print(f"Erro ao ler {file_path}: {e}")

    return mappings

def send_bulk_mappings(mappings):
    if not mappings:
        return 0

    data = json.dumps({"mappings": mappings}).encode("utf-8")
    req = urllib.request.Request(API_URL, data=data, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req) as response:
            res = json.loads(response.read().decode("utf-8"))
            return res.get("processed", 0)
    except Exception as e:
        print(f"Erro ao enviar para API: {e}")
        return 0

def main():
    print(f"Iniciando Zomboid Workshop Scanner (Workers: {MAX_WORKERS}, API: {API_URL})")
    # Exemplo de lote para teste/processamento
    print("Scanner pronto para receber comandos ou executar varredura de diretório.")

if __name__ == "__main__":
    main()
