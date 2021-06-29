import json
import os
import pika
import unittest
import uuid

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

class RPCClient(object):

    def __init__(self, rpc_queue):
        self.rpc_queue = rpc_queue

        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host='localhost'))

        self.channel = self.connection.channel()

        result = self.channel.queue_declare(queue='', exclusive=True)
        self.callback_queue = result.method.queue

        self.channel.basic_consume(
            queue=self.callback_queue,
            on_message_callback=self._on_response,
            auto_ack=True)

    def close(self):
        self.channel.close()
        self.connection.close()

    def _on_response(self, ch, method, props, body):
        if self.corr_id == props.correlation_id:
            self.response = body

    def call(self, payload):
        self.response = None
        self.corr_id = str(uuid.uuid4())
        self.channel.basic_publish(
            exchange='',
            routing_key=self.rpc_queue,
            properties=pika.BasicProperties(
                reply_to=self.callback_queue,
                correlation_id=self.corr_id,
            ),
            body=json.dumps(payload))
        while self.response is None:
            self.connection.process_data_events()
        return self.response

    def info(self):
        return self.call({"cmd": "info", "params": {}})

    def ping(self):
        return self.call({"cmd": "ping", "params": {}})


class AudioMetadataReader(RPCClient):
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
        dir = os.path.abspath(f"file/testdata/{ext}")
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
        self.assertEqual(list(r["actors"].keys())[0], "test_performer")
        self.assertEqual(r["discs"][0]["number"], 1)
        self.assertEqual(r["total_tracks"], 10)
        self.assertEqual(r["tracks"][0]["position"], "03")
        self.assertEqual(list(t["composition"]["actor_roles"].keys())[0], "test_composer")
        self.assertEqual(list(t["record"]["actors"].keys())[0], "test_track_artist")
        self.assertEqual(t["record"]["genres"][0], "test_genre")
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
            self.assertEqual(r["pictures"][0]["pict_type"], "cover_front")
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
