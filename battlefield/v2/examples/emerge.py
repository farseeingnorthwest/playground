import subprocess
import sys
import json
import requests
import click

class Emerge:
    def __init__(self, mod, url):
        self.mod = mod if mod is not None else "github.com/farseeingnorthwest/playground/battlefield/v2/examples/json"
        self.url = url

    def serve(self, name, sect, *index):
        process = subprocess.run(
            ["go", "run", self.mod, sect, *[str(i) for i in index]],
            capture_output=True,
            text=True,
        )
        if process.returncode != 0:
            return

        reactor = json.loads(process.stdout)
        skill = dict(name=name, reactor=reactor)
        if self.url is None:
            json.dump(skill, sys.stdout)
            print("\x1e")
        else:
            r = requests.post(self.url, json=skill)
            print(r.text)


@click.command()
@click.option('--mod', type=click.Path(exists=True))
@click.option('--url')
def main(mod, url):
    e = Emerge(mod, url)
    for (name, i) in [
            ("Normal Attack", 0),
            ("Critical / Set", 1),
            ("Critical / Buff", 2),
            ("Element Theory", 3),
    ]:
        e.serve(name, "regular", i)

    for i, name in enumerate([
            "織田",
            "豐臣",
            "上杉",
            "徳川",
            "武田",
            "梅花",
            "鑽石",
            "王牌",
            "紅心",
            "黑桃",
    ]):
        for j in range(0, 4):
            e.serve(f"{name}-{1 + j}", "special", i, j)


if __name__ == "__main__":
    main()
