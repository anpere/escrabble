
import gizeh

white = (1, 1, 1)
black = (0, 0, 0)
class Board(object):
    filename = "genBoard1.png"

    def gen_piece_generator(self, piece_count):
        for character, count in piece_count.items():
            for i in range(count):
                yield Piece(character, count)

    def getNextPiece(self):
        return self.piece_generator.__next__()

    def __init__(self, piece_map, piece_count, width, height):
        self.pieces = []
        self.width = width
        self.height = height
        self.piece_generator = self.gen_piece_generator(piece_count)

    def draw(self, surface):
        for piece in self.pieces:
            piece.draw(surface)

    def gen_board(self):
        self.pieces = []
        num_pieces = 0
        self.pieces = []
        rows = int(self.height / Piece.height)
        cols = int(self.width / Piece.width)
        for col in range(cols):
            for row in range(rows):
                try:
                    piece = self.getNextPiece()
                    piece.set_draw_properties(
                        xy=[piece.height/2 + col*piece.height, piece.width/2 + row*piece.width],
                        )
                    self.pieces.append(piece)
                except:
                    break

        surface = gizeh.Surface(width=int(cols*Piece.width), height=int(rows*Piece.height))
        self.draw(surface)
        surface.get_npimage()
        surface.write_to_png(self.filename)
        print("did the generation")


class Piece(object):

    fontfamily="Impact"
    letterColor=black
    rectangleColor=white
    height=500
    width=height
    lettersize=height/3
    numbersize=lettersize/3

    def __init__(
        self,
        character,
        value,
        **kwargs
    ):
        self.character = character
        self.value = value
        for key, value in kwargs.items():
            self.__setattr__(key, value)

    def set_draw_properties(self, *args, **kwargs):
        for attr, val in kwargs.items():
            self.__setattr__(attr, val)

    def draw(self, surface):
        """
        Draws a piece onto a surface
        """
        border = gizeh.rectangle(
            lx=self.height,
            ly=self.height,
            xy=self.xy,
            fill=black
        )
        inner = gizeh.rectangle(
            lx=self.width*.95,
            ly=self.height*.95,
            xy=self.xy,
            fill=self.rectangleColor,
        )
        letter = gizeh.text(
            self.character,
            fontfamily=self.fontfamily,
            fontsize=self.lettersize,
            fill=self.letterColor,
            xy=(self.xy[0], self.xy[1]),
            )
        value = gizeh.text(
            str(self.value),
            fontfamily=self.fontfamily,
            fontsize=self.numbersize,
            fill=self.letterColor,
            xy=(self.xy[0]+self.width*.4, self.xy[1]+self.height*.4),
        )
        figures = [border, inner, letter, value]
        for figure in figures:
            figure.draw(surface)

def pq(num):
    factors = [divisor for divisor in range(1, num) if num%divisor ==0]
    p = factors[int(len(factors)/2)]
    return p, num/p

def get_pieces(file_name):
    """
    File name of the piece data
    """
    f = open(file_name, 'r')
    piece_map = {}
    for line in f:
        character, value = line.split(":")
        piece_map[character] = int(value)
    return piece_map

if __name__ == "__main__":
    pieces = get_pieces("pieces.json")

    if True:
        piece_map = get_pieces("piece_value.txt")
        piece_count = get_pieces("piece_count.txt")
    else:
        piece_map = {
            "a":"10",
            "b":"8",
            "c":"4",
        }
        piece_count = {
            "a":4,
            "b":3,
            "c":1,
        }
    num_pieces = 0
    for piece, count in piece_count.items():
        num_pieces += count

    rows, cols = pq(num_pieces)
    # board.gen_board()

    a = Piece("a", "1", xy=[Piece.width*.5, Piece.height*.5])
    surface = gizeh.Surface(width=a.width, height=a.height)
    a.draw(surface)
    surface.get_npimage()
    surface.write_to_png("a.png")

    board = Board(piece_map, piece_count, width=Piece.width*rows, height=Piece.height*cols)
    board.gen_board()
