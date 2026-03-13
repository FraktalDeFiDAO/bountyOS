import subprocess
import time
from flask import Flask, render_template, jsonify
from pymongo import MongoClient

app = Flask(__name__)
client = MongoClient("mongodb://localhost:27017/")
db = client.twitter_trends

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/fetch-trends', methods=['GET'])
def fetch_trends():
    
    subprocess.run(['python', 'selenium_script.py'], check=True)  

    
    time.sleep(5)

    latest = db.trends.find_one(sort=[("date_time", -1)])

    return jsonify({
        "trends": latest["trends"],
        "date_time": latest["date_time"].strftime('%Y-%m-%d %H:%M:%S'),
        "ip_address": latest["ip_address"]
    })

if __name__ == '__main__':
    app.run(debug=True)
