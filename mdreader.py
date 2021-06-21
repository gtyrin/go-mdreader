import common
import json
import os
import unittest

# Данные для сверки значений тегов для аудиофайлов 440_hz_mono.*:
# 					                mp3	    flac	wv	    dsf	    значение
# - title					        +	    +	    +		        test_album_title
# - work
# 	- composer			            +	    +	    +		        test_composer
# - recording
# 	- performer			            +	    +	    +		        test_performer
# 	- genres			            +	    +	    +		        test_genre
# - totaltracks			            + 	    +	    +		        10
# - tracks
# 	- position			            +	    +	    +		        3
# 	- title				            +		+	    +	            test_track_title
# 	- artist			            +       +		+		        test_track_artist
# - country				            +	    +	    +		        test_country
# - label					        +	    +	    +		        test_label
# - catno					        +	    +	    +		        test_catno
# - year					        +	    +	    +		        2000
# - notes					        +	    +	    +		        test_notes
# - rutracker, DISCOGS_RELEASE_ID	+	    +	    +		        123456789
# - cover					        +	    +	    +


class AudioMetadataReader(common.RPCClient):
    def __init__(self):
        super().__init__('mdreader')

    def release(self, dir):
        return self.call({"cmd": "release", "params": {"dir": dir}})


class TestMetadataReader(unittest.TestCase):
    def setUp(self):
        self.r = AudioMetadataReader()

    def tearDown(self):
        self.r.close()

    def test_ping(self):
        self.assertEqual(self.r.ping(), b'')

    def release(self, ext):
        dir = f"testdata/audio/file/{ext}"
        self.assertTrue(os.path.exists(dir))
        resp = json.loads(self.r.release(dir))
        if resp.get("error"):
            raise BaseException(resp)
        else:
            return resp

    def check_resp(self, resp):
        r = resp["release"]
        t = resp["release"]["tracks"][0]
        self.assertEqual(r["title"], "test_album_title")
        self.assertEqual(r["work"]["actors"][0]["name"], "test_composer")
        self.assertEqual(r["recording"]["actors"][0]["name"], "test_performer")
        self.assertEqual(r["recording"]["genres"][0], "test_genre")
        self.assertEqual(r["discs"][0]["number"], 1)
        self.assertEqual(r["total_tracks"], 10)
        self.assertEqual(t["record"]["actors"][0]["name"], "test_track_artist")
        self.assertEqual(resp["release"]["tracks"][0]["position"], "03")
        self.assertEqual(t["title"], "test_track_title")
        basename, ext = (os.path.splitext(t["file_info"]["file_name"]))
        if ext not in (".mp3", ".wv"):  # TODO
            self.assertEqual(t["duration"], 500)
        self.assertEqual(basename, "440_hz_mono")
        if ext not in (".dsf", ".wv"):  # TODO
            self.assertEqual(t["audio_info"]["samplerate"], 44100)
            self.assertEqual(t["audio_info"]["sample_size"], 16)
        if ext != ".wv":  # TODO
            self.assertEqual(t["audio_info"]["channels"], 1)
            self.assertEqual(
                r["pictures"][0]["pict_meta"]["mime_type"],
                "image/jpeg")
            self.assertEqual(r["pictures"][0]["pict_type"], 3)
        self.assertEqual(r["publishing"][0]["name"], "test_label")
        self.assertEqual(r["publishing"][0]["catno"], "test_catno")
        self.assertEqual(r["country"], "test_country")
        self.assertEqual(r["year"], 2000)
        self.assertEqual(r["notes"], "test_notes")

    def test_mp3(self):
        self.check_resp(self.release("mp3"))

    def test_flac(self):
        self.check_resp(self.release("flac"))

    def test_dsf(self):
        self.check_resp(self.release("dsf"))

    def test_wavpack(self):
        self.check_resp(self.release("wavpack"))


if __name__ == '__main__':
    unittest.main()
