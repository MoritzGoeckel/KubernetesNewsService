import os
from pymongo import MongoClient
from pymongo.errors import ServerSelectionTimeoutError
from NgramLanguageModel import Model

def main():
    print('Language Model version 0.01')

    mongoClient = getConnection()
    if mongoClient:
        collection_articles = mongoClient.news.articles

        latest = 0
        print('latest: ', latest)

        articles = collection_articles.find({'article_perplexity': {'$exists': False}})
        model = Model(n=3)
        model.read_frequencies(path='frequencies/reuters_adjusted_freq') #TODO: should this path be part of the env-variables?
        c = 0 #needs to be done this way instead of len(list(articles)) to avoid memory problems
        for article in articles:
            article_content = article['content']
            if article_content:
                article_date_time = article['datetime']
                article_url = article['url']
                article_id = article['_id']
                pp = model.perplexity(article_content)
                collection_articles.update_one({'_id':article_id}, {'$set': {'article_perplexity': pp}})
            c += 1
        print('Processed', c, 'articles')


def getConnection():
    mongoUrl = os.environ.get('mongo-url')
    if not mongoUrl:
        print('Environment variables are not set')
    print('mongo url: ', mongoUrl)

    mongoPw = os.environ.get('mongo-pw')
    mongoUser = os.environ.get('mongo-user')
    print('mongo credentials: ', mongoUser, mongoPw)
    print('Connecting to mongo...')

    mp = 27017 #default mongo port
    client = MongoClient(host=mongoUrl, port=mp, username=mongoUser, password=mongoPw)

    try:
        client.server_info()
    except ServerSelectionTimeoutError as err:
        print(err)
        return

    print('Mongo connection established')
    return client


if __name__ == "__main__":
    main()