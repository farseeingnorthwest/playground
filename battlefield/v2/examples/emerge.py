import subprocess
import json
import requests
import click

class Emerge:
    def __init__(self, src, dst):
        self.src = src
        self.dst = dst

    def serve(self, name, sect, *index):
        process = subprocess.run(
            ["go", "run", self.src, sect, *[str(i) for i in index]],
            capture_output=True,
            text=True,
        )
        if process.returncode != 0:
            return

        reactor = json.loads(process.stdout)
        r = requests.post(self.dst, json=dict(name=name, reactor=reactor))
        print(r.text)


@click.command()
@click.argument('src', type=click.Path(exists=True))
@click.argument('dst')
def main(src, dst):
    e = Emerge(src, dst)
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
