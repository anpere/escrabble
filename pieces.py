
import gizeh
import cairocffi as cairo
import numpy

white = (1, 1, 1)
black = (0, 0, 0)
red   = (1, 0, 0)
class Apl(object):
    def __init__(
        self,
        sizex,
        sizey,
        **kwargs
    ):
        # load png and convert to numpy array

        image_surface = cairo.ImageSurface.create_from_png("logo.png")
        im = 0+numpy.frombuffer(image_surface.get_data(), numpy.uint8)
        im.shape = (image_surface.get_height(), image_surface.get_width(), 4)
        im = im[:,:,[2,1,0,3]] # put RGB back in order
        gizeh_pattern = gizeh.ImagePattern(im)
        self.height = im.shape[0]
        self.width = im.shape[1]
        gizeh_pattern.scale(rx=sizex, ry=sizey)
        self.scalex = sizex
        self.scaley = sizey
        self.im = im
        self.gizeh_pattern = gizeh_pattern
        for key, value in kwargs.items():
            self.__setattr__(key, value)

    def draw(self, surface):

        r = gizeh.rectangle(lx=self.width*10, ly=self.height*10, xy=self.xy, fill=self.gizeh_pattern)
        r.scale(rx=self.scalex, ry=self.scaley)
        r.draw(surface)

class Board(object):
    filename = "board.png"

    @property
    def pieces(self):
        characters = sorted(self.piece_points.keys())
        for character in characters:
            print(character)
            for i in range(self.piece_count[character]):
                yield Piece(character, self.piece_points[character])

    def getNextPiece(self):
        return next(self.pieces)

    def __init__(self, piece_points, piece_count, width, height):
        self.width = int(width)
        self.height = int(height)
        self.piece_points = piece_points
        self.piece_count = piece_count

    def gen_board(self):
        rows = int(self.height / Piece.height)
        cols = int(self.width / Piece.width)
        print("generating board {} by {}".format(rows, cols))
        surface = gizeh.Surface(width=self.width, height=self.height)
        back_surface = gizeh.Surface(width=self.width, height=self.height)
        i = 0
        for piece in self.pieces:
            row = int(i/cols)
            col = i%cols
            x_y=[piece.width/2 + col*piece.width, piece.height/2 + row*piece.height]
            piece.set_draw_properties(
                xy=x_y,
            )
            piece.draw(surface)
            i+=1
        for row in range(rows+1):
            for col in range(2*cols+1):
                tiny_logo = Apl(2,2)
                tiny_logo.xy = [tiny_logo.width * row, tiny_logo.height  * col]
                tiny_logo.draw(back_surface)
        surface.get_npimage()
        surface.write_to_png(self.filename)
        back_surface.get_npimage()
        back_surface.write_to_png("back_board.png")
        print("did the generation")



class Piece(object):

    fontfamily="Impact"
    letterColor=black
    rectangleColor=white
    height=200
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
            fill=red
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
    print("pieces.py...")
    if True:
        piece_map = get_pieces("piece_value.txt")
        piece_count = get_pieces("piece_count.txt")
    else:
        piece_map = {"a":"10","b":"8","c":"4",}
        piece_count = {"a":4,"b":3,"c":1,}

    print("pieces: {}".format(piece_count))
    num_pieces = 0
    for piece, count in piece_count.items():
        num_pieces += count
    rows, cols = pq(num_pieces)
    if rows == 1:
        # prime numbers, fuck that shit.
        rows,cols = pq(num_pieces+3)
    print("Num pieces {} generating board {} rows by {} columns".format(num_pieces, rows, cols))


    board = Board(piece_map, piece_count, width=Piece.width*cols, height=Piece.height*rows)
    board.gen_board()
    logo = Apl(1,4,xy=(0,0))
    surface = gizeh.Surface(logo.width, logo.height,  bg_color=(0, 0.0, 0.0))
    logo.draw(surface)
