import random

import numpy as np
import torch
from torch.nn import CrossEntropyLoss
from torchvision import datasets, transforms
from torchvision.models.mobilenet import mobilenet_v2


def tst(model, device, test_loader):
    model.eval()
    test_loss = 0
    correct = 0
    loss_func = CrossEntropyLoss()
    with torch.no_grad():
        for data, target in test_loader:
            data, target = data.to(device), target.to(device)
            output = model(data)
            test_loss += loss_func(output, target)
            pred = output.argmax(dim=1, keepdim=True)  # get the index of the max log-probability
            correct += pred.eq(target.view_as(pred)).sum().item()

    test_loss /= len(test_loader)

    print('\nTest set: Average loss: {:.4f}, Accuracy: {}/{} ({:.0f}%)\n'.format(
        test_loss, correct, len(test_loader.dataset),
        100. * correct / len(test_loader.dataset)))


def main():
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

    kwargs = {}
    if torch.cuda.is_available():
        kwargs.update({'num_workers': 1, 'pin_memory': True})

    transform = transforms.Compose([
        transforms.ToTensor(),
        transforms.Normalize((0.1307,), (0.3081,))
    ])

    dataset2 = datasets.CIFAR10('/tmp/cifar10', train=False, transform=transform, download=True)
    
    test_loader = torch.utils.data.DataLoader(dataset2, shuffle=False, **kwargs)

    model = mobilenet_v2(pretrained=True)
    model.classifier[1] = torch.nn.Linear(in_features=model.classifier[1].in_features, out_features=10)
    model.to(device)

    tst(model, device, test_loader)

if __name__ == '__main__':
    SEED=0
    random.seed(SEED)     # python random generator
    np.random.seed(SEED)  # numpy random generator

    torch.manual_seed(SEED)
    torch.use_deterministic_algorithms(True)

    main()
    
