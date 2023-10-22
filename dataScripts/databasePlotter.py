from matplotlib import pyplot as plt
import sys

def readDataFromFile(fileName):
    with open(fileName, "r") as file:
        lines = file.readlines()
        header = lines.pop(0)
        
        data = []
        for line in lines:
            data.append(line.split())
    return data

def CreateProductEntryDictionary(data):
    productDictionary = dict()

    for product in data:
        productName = product[0]
        productAmount = int(product[1])
        productPrice = float(product[2])
    
        entry = (productAmount, productPrice)
    
        if productName in productDictionary:
            productDictionary[productName].append(entry)
        else:
            productDictionary[productName] = [entry]
            
    return productDictionary

def plotProduct(productName, entry):
    s = 1
    prices = [price for _, price in entry]
    amount = [amount for amount, _ in entry]
    entriesIndex = range(1, len(prices)+1)
    
    plt.scatter(entriesIndex, prices, marker='.', s=s)
    plt.title(f"Prices change for {productName} over time")
    plt.yscale("log")
    plt.savefig(f"plots/{productName}_prices_change")
    
    plt.close()
    
    plt.scatter(entriesIndex, amount, marker='.', s=s)
    plt.title(f"Amount change for {productName} over time")
    plt.yscale("log")
    plt.savefig(f"plots/{productName}_amount_change")



if __name__ == "__main__":
    # filename = sys.argv[1]
    fileName = "data/exampleDatabase"
    data = readDataFromFile(fileName)        

    productDictionary = CreateProductEntryDictionary(data)

    for productName, entry in productDictionary.items():
        plotProduct(productName, entry)