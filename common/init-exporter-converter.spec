################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Utility for converting init-exporter procfiles from v1 to v2 format
Name:           init-exporter-converter
Version:        0.12.0
Release:        0%{?dist}
Group:          Development/Tools
License:        MIT
URL:            https://github.com/funbox/init-exporter-converter

Source0:        %{name}-%{version}.tar.gz

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.19

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
Utility for exporting services described by Procfile to init system.

################################################################################

%prep
%setup -q

%build
if [[ ! -d "%{name}/vendor" ]] ; then
  echo "This package requires vendored dependencies"
  exit 1
fi

pushd %{name}
  %{__make} %{?_smp_mflags} all
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -pm 755 %{name}/%{name} %{buildroot}%{_bindir}/

%clean
rm -rf %{buildroot}

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_bindir}/init-exporter-converter

################################################################################

%changelog
* Fri Mar 10 2023 Anton Novojilov <andyone@fun-box.ru> - 0.12.0-0
- Added verbose version output
- Dependencies update
- Code refactoring

* Fri Apr 01 2022 Anton Novojilov <andyone@fun-box.ru> - 0.11.2-0
- Removed pkg.re usage
- Added module info
- Added Dependabot configuration

* Mon Jan 10 2022 Anton Novojilov <andyone@fun-box.ru> - 0.11.1-0
- Minor UI improvements

* Thu Jul 30 2020 Anton Novojilov <andyone@fun-box.ru> - 0.11.0-0
- Package ek updated to v12
- Package go-simpleyaml updated to v2

* Mon Aug 05 2019 Anton Novojilov <andyone@fun-box.ru> - 0.10.1-0
- Updated for compatibility with the latest version of ek package

* Thu Jan 10 2019 Anton Novojilov <andyone@fun-box.ru> - 0.10.0-0
- Package ek updated to v10

* Thu Nov 01 2018 Anton Novojilov <andyone@fun-box.ru> - 0.9.1-0
- Rebuilt with the latest procfile parser

* Thu Apr 05 2018 Anton Novojilov <andyone@fun-box.ru> - 0.9.0-1
- Rebuilt with Go 1.10

* Tue Sep 19 2017 Anton Novojilov <andyone@fun-box.ru> - 0.9.0-0
- Improved environment variables parsing in v1

* Fri May 19 2017 Anton Novojilov <andyone@fun-box.ru> - 0.8.0-0
- Migrated to ek.v9
- Added support of multiple file converting
